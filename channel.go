package evil

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/gorilla/websocket"
)

type JoinMessage struct {
	// Flash
	// params
	Session string `json:"session"`
	Static  string `json:"static"`
	URL     string `json:"url"`
}

type Channel struct {
	Conn           *Conn
	Renderer       *Renderer
	WebSocket      *websocket.Conn
	SessionFactory *SessionFactory

	sessions sync.Map
}

func (c *Channel) Run() error {
	for {
		var msg Message
		if err := c.WebSocket.ReadJSON(&msg); err != nil {
			return errors.Wrap(err, "unmarshaling message")
		}

		resp, err := c.handleMessage(&msg)
		if err != nil {
			return errors.Wrapf(err, "handling message %#v", msg)
		}

		if err := c.WebSocket.WriteJSON(resp); err != nil {
			return errors.Wrap(err, "writing response")
		}
	}
}

func (c *Channel) handleMessage(msg *Message) (*Message, error) {
	switch msg.Event {
	case "phx_join":
		return c.handleJoin(msg)
	case "event":
		return c.handleEvent(msg)
	case "heartbeat":
		return c.handleHeartbeat(msg)
	default:
		log.Printf("unhandled message: %#v", msg)
		return msg, nil
	}
}

func (c *Channel) handleHeartbeat(msg *Message) (*Message, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"status":   "ok",
		"response": map[string]interface{}{},
	})

	if err != nil {
		return nil, err
	}

	return &Message{
		JoinRef: nil,
		Ref:     msg.Ref,
		Topic:   msg.Topic,
		Event:   "phx_reply",
		Payload: payload,
	}, nil
}

func (c *Channel) handleJoin(msg *Message) (*Message, error) {
	var evt JoinMessage
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return nil, err
	}

	s, err := DecodeSession(evt.Session)
	if err != nil {
		return nil, err
	}

	fmt.Printf("JOIN: %#v\n", evt)

	// Need to distinquish between session sig and actual session but w/e
	s, err = c.SessionFactory.LoadSession(s.ID)
	if err != nil {
		return nil, err
	}

	if err := c.Renderer._views[s.View].OnMount(s); err != nil {
		return nil, err
	}

	out, err := c.Renderer.RenderView(s.View, s)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(map[string]interface{}{
		"status": "ok",
		"response": map[string]interface{}{
			"rendered": &Diff{
				Dynamic: []string{out},
				Static:  []string{"", ""},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &Message{
		JoinRef: msg.JoinRef,
		Ref:     msg.Ref,
		Topic:   msg.Topic,
		Event:   "phx_reply",
		Payload: payload,
	}, nil
}

func (c *Channel) handleEvent(msg *Message) (*Message, error) {
	var evt Event
	if err := json.Unmarshal(msg.Payload, &evt); err != nil {
		return nil, errors.Wrap(err, "unmarshaling event")
	}

	session, err := c.SessionFactory.LoadSession(msg.Topic[3:])
	if err != nil {
		return nil, err
	}

	if err := c.Renderer._views[session.View].HandleEvent(session, &evt); err != nil {
		return nil, err
	}

	out, err := c.Renderer.RenderView(session.View, session)
	if err != nil {
		return nil, errors.Wrap(err, "rendering")
	}

	payload, err := json.Marshal(map[string]interface{}{
		"status": "ok",
		"response": map[string]interface{}{
			"diff": &Diff{
				Dynamic: []string{out},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return &Message{
		JoinRef: msg.JoinRef,
		Ref:     msg.Ref,
		Topic:   msg.Topic,
		Event:   "phx_reply",
		Payload: payload,
	}, nil
}
