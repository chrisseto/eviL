package evil

import (
	"encoding/json"

	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
	"github.com/cockroachdb/errors"
)

type LiveViewChannel struct {
	Templates      map[string]*template.Template
	SessionFactory *SessionFactory
}

var _ channel.Channel = &LiveViewChannel{}

func (c *LiveViewChannel) Join(_ *channel.Session, m *channel.Message) (interface{}, error) {
	var j channel.Join
	if err := json.Unmarshal(m.Payload, &j); err != nil {
		return nil, err
	}

	// TODO should actually validate this
	// TODO mix with melody sessions
	session, err := c.SessionFactory.FromToken(j.Session)
	if err != nil {
		return nil, err
	}

	// if err := c.Renderer.Mount(session.View, session); err != nil {
	// 	return nil, err
	// }

	template, ok := c.Templates[session.View]
	if !ok {
		return nil, errors.Newf("no such view: %s", session.View)
	}

	diff := template.Render(session)

	return diff, nil
}

func (c *LiveViewChannel) Handle(_ *channel.Session, m *channel.Message) (interface{}, error) {
	var e channel.Event
	if err := json.Unmarshal(m.Payload, &e); err != nil {
		return nil, err
	}

	session, err := c.SessionFactory.LoadSession(m.Topic[3:])
	if err != nil {
		return nil, err
	}

	// if err := c.Renderer.Event(s.View, s, &e); err != nil {
	// 	return nil, err
	// }

	template, ok := c.Templates[session.View]
	if !ok {
		return nil, errors.Newf("no such view: %s", session.View)
	}

	diff := template.Render(session)

	return diff, nil
}
