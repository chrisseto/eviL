package evil

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// TODO figure out binary encoding
type Message struct {
	JoinRef *string
	Ref     string
	Topic   string
	Event   string
	Payload json.RawMessage
}

var _ json.Marshaler = &Message{}

func (m *Message) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&m.JoinRef, &m.Ref, &m.Topic, &m.Event, &m.Payload}
	// wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	// if g, e := len(tmp), wantLen; g != e {
	// 	return fmt.Errorf("wrong number of fields in Notification: %d != %d", g, e)
	// }
	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		m.JoinRef,
		m.Ref,
		m.Topic,
		m.Event,
		m.Payload,
	})
}

type Diff struct {
	// Seems like there is a reply "r"
	// and a events "e" list as well? haven't seen them used yet.
	Dynamic []string
	Static  []string
}

func (d *Diff) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")
	enc := json.NewEncoder(&buf)

	for i, dynamic := range d.Dynamic {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf(`"%d":`, i))
		if err := enc.Encode(dynamic); err != nil {
			return nil, err
		}
	}

	if d.Static != nil {
		if len(d.Dynamic) > 0 {
			buf.WriteString(`, `)
		}

		buf.WriteString(`"s":`)

		if err := enc.Encode(d.Static); err != nil {
			return nil, err
		}
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}
