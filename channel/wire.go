package channel

import (
	"encoding/json"
)

type EventType string
type Status string

const (
	StatusOK    Status = "ok"
	StatusError Status = "error"

	TypeHearbeat EventType = "heartbeat"
	TypeEvent    EventType = "event"
	TypeJoin     EventType = "phx_join"
	TypeReply    EventType = "phx_reply"
	TypeLeave    EventType = "phx_leave"
)

type Message struct {
	JoinRef *string
	Ref     *string
	Topic   string
	Event   EventType
	Payload json.RawMessage
}

var _ json.Marshaler = &Message{}
var _ json.Unmarshaler = &Message{}

func (m *Message) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&m.JoinRef, &m.Ref, &m.Topic, &m.Event, &m.Payload}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}

	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		&m.JoinRef,
		&m.Ref,
		&m.Topic,
		&m.Event,
		&m.Payload,
	})
}

// TODO combine with Message
type Reply struct {
	JoinRef *string     `json:"join_ref"`
	Ref     *string     `json:"ref"`
	Topic   string      `json:"topic"`
	Status  Status      `json:"status"`
	Payload interface{} `json:"payload"`
}

func (r *Reply) SetPayload(resp interface{}, err error) {
	if err != nil {
		r.Status = StatusError
		r.Payload = err // Might need some nice handling here
	} else {
		r.Status = StatusOK
		r.Payload = resp
	}
}

func (r *Reply) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		&r.JoinRef,
		&r.Ref,
		&r.Topic,
		"phx_reply",
		// TODO Ommit if empty ??
		map[string]interface{}{
			"status":   r.Status,
			"response": r.Payload,
		},
	})
}

type Response struct {
	Status   Status      `json:"status"`
	Response interface{} `json:"response"`
}

type Join struct {
	URL     string                 `json:"url"`
	Params  map[string]interface{} `json:"params"`
	Session string                 `json:"session"`
	Static  string                 `json:"static"`
}

type Event struct {
	Type  string `json:"type"`
	Event string `json:"event"`
	Value string `json:"value"`
}
