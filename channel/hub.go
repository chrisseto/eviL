package channel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/olahol/melody"
)

type Hub struct {
	router *melody.Melody

	mu       sync.Mutex
	channels map[string]Channel
}

var _ http.Handler = &Hub{}

func NewHub() *Hub {
	h := &Hub{
		router:   melody.New(),
		channels: make(map[string]Channel),
	}

	h.router.Config.MaxMessageSize = 1024 * 1024

	h.router.HandleMessage(func(s *melody.Session, data []byte) {
		if err := h.onMessage(s, data); err != nil {
			panic(err)
		}
	})

	h.router.HandleClose(func(s *melody.Session, _ int, _ string) error {
		if err := h.onLeave(s); err != nil {
			panic(err)
		}
		return nil
	})

	return h
}

func (h *Hub) Register(pattern string, channel Channel) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.channels[pattern[:len(pattern)-2]] = channel
}

func (h *Hub) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := h.router.HandleRequest(rw, r); err != nil {
		panic(err)
	}
}

func (h *Hub) Broadcast(pattern string, event string, payload interface{}) error {
	marshalled, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	// This won't work with globbing :/
	msg := Message{
		// Hack
		Topic:   "lv:" + pattern,
		Event:   EventType(event),
		Payload: marshalled,
	}

	data, err := json.Marshal(&msg)
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	h.router.BroadcastFilter(data, func(session *melody.Session) bool {
		// TODO allow globbing
		return session.MustGet("id") == pattern
	})

	return nil
}

func (h *Hub) onMessage(session *melody.Session, data []byte) error {
	var e Message
	if err := json.Unmarshal(data, &e); err != nil {
		return errors.Wrap(err, "unmarshaling message")
	}

	fmt.Printf("recieved message: %#v\n", e)

	reply := &Reply{
		JoinRef: e.JoinRef,
		Ref:     e.Ref,
		Topic:   e.Topic,
	}

	var err error
	var resp interface{}

	if e.Topic != "phoenix" {
		// TODO refactor into a matching function
		topicName := strings.SplitN(e.Topic, ":", 2)[0]
		channel, ok := h.channels[topicName]
		if !ok {
			fmt.Printf("channels: %#v\n", h.channels)
			return errors.Newf("no such topic: %s", e.Topic)
		}

		switch e.Event {
		case TypeJoin:
			resp, err = channel.Join(session, &e)
		default:
			resp, err = channel.Handle(session, &e)
		}
	} else {
		// Giant switch might be nicer
		switch e.Event {
		case TypeHearbeat:
			// Should return an empty reply
			resp = map[string]interface{}{}
		default:
			return errors.Newf("unknown phoenix event: %s", e.Event)
		}
	}

	if resp == nil && err == nil {
		return nil
	}

	if err != nil {
		// TODO upgrade to logging
		fmt.Printf("Failed to handle message: %+v\n", err)
	}

	reply.SetPayload(resp, err)

	if err := writeJSON(session, reply); err != nil {
		return errors.Wrap(err, "writing JSON")
	}

	return nil
}

func (h *Hub) onLeave(session *melody.Session) error {
	// TODO
	return nil
}
