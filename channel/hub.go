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

	h.router.HandleMessage(func(s *melody.Session, data []byte) {
		if err := h.onMessage(s, data); err != nil {
			panic(err)
		}
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

func (h *Hub) onMessage(session *melody.Session, data []byte) error {
	var e Message
	if err := json.Unmarshal(data, &e); err != nil {
		return errors.Wrap(err, "unmarshaling message")
	}

	fmt.Printf("Recieved Message: %#v\n", e)

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
		case TypeHearbeat:
			// Should validate topic == phoenix ??
			// Should return an empty reply
		default:
			resp, err = channel.Handle(session, &e)
		}
	} else {
		switch e.Event {
		case TypeHearbeat:
			resp = map[string]interface{}{}
		default:
			return errors.Newf("unknown phoenix event: %s", e.Event)
		}
	}

	if resp == nil && err == nil {
		return nil
	}

	if err != nil {
		reply.Status = StatusError
		reply.Payload = err
		fmt.Printf("Failed to handle message: %#v\n", e)
	} else {
		reply.Status = StatusOK
		reply.Payload = resp
	}

	if err := writeJSON(session, reply); err != nil {
		return err
	}

	return nil
}
