package evil

import (
	"encoding/json"

	"github.com/chrisseto/evil/channel"
	"github.com/olahol/melody"
)

type LiveViewChannel struct {
	Renderer       *Renderer
	SessionFactory *SessionFactory
}

var _ channel.Channel = &LiveViewChannel{}

func (c *LiveViewChannel) Join(_ *melody.Session, m *channel.Message) (interface{}, error) {
	var j channel.Join
	if err := json.Unmarshal(m.Payload, &j); err != nil {
		return nil, err
	}

	// TODO should actually validate this
	session, err := c.SessionFactory.FromToken(j.Session)
	if err != nil {
		return nil, err
	}

	if err := c.Renderer.Mount(session.View, session); err != nil {
		return nil, err
	}

	out, err := c.Renderer.RenderView(session.View, session)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"rendered": &Diff{
			Dynamic: []string{out, ""},
			Static:  []string{" ", " "},
		},
	}, nil
}

func (c *LiveViewChannel) Handle(_ *melody.Session, m *channel.Message) (interface{}, error) {
	var e channel.Event
	if err := json.Unmarshal(m.Payload, &e); err != nil {
		return nil, err
	}

	s, err := c.SessionFactory.LoadSession(m.Topic[3:])
	if err != nil {
		return nil, err
	}

	if err := c.Renderer.Event(s.View, s, &e); err != nil {
		return nil, err
	}

	out, err := c.Renderer.RenderView(s.View, s)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"diff": &Diff{
			Dynamic: []string{out, ""},
			Static:  []string{" ", " "},
		},
	}, nil
}
