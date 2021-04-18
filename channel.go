package evil

import (
	"encoding/json"

	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
	"github.com/cockroachdb/errors"
	"github.com/dgrijalva/jwt-go"
)

type LiveViewChannel struct {
	Template       *template.Template
	SessionFactory *SessionFactory
	Views          map[string]View
}

var _ channel.Channel = &LiveViewChannel{}

var secret = []byte("password123")

func (c *LiveViewChannel) RegisterView(name string, view View) {
	c.Views[name] = view
}

func (c *LiveViewChannel) verifySession(signed string) (*SessionClaims, error) {
	token, err := jwt.ParseWithClaims(signed, &SessionClaims{}, func(t *jwt.Token) (interface{}, error) {
		_ = t.Method.(*jwt.SigningMethodHMAC)
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*SessionClaims)

	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}

func (c *LiveViewChannel) Join(session *channel.Session, m *channel.Message) (interface{}, error) {
	var j channel.Join
	if err := json.Unmarshal(m.Payload, &j); err != nil {
		return nil, err
	}

	claims, err := c.verifySession(j.Session)
	if err != nil {
		return nil, err
	}

	if view, ok := c.Views[claims.View]; ok {
		if err := view.OnMount(session); err != nil {
			return nil, err
		}
	}

	diff, err := c.Template.ExecuteTemplate(claims.View, session.Keys())
	if err != nil {
		return nil, errors.Wrap(err, "executing template")
	}

	return map[string]interface{}{
		"rendered": diff,
	}, nil
}

func (c *LiveViewChannel) Handle(session *channel.Session, m *channel.Message) (interface{}, error) {
	var e channel.Event
	if err := json.Unmarshal(m.Payload, &e); err != nil {
		return nil, err
	}

	// if view, ok := c.Views[claims.View]; ok {
	// 	if err := view.OnMount(session); err != nil {
	// 		return nil, err
	// 	}
	// }

	diff, err := c.Template.ExecuteTemplate("", session)

	return diff, errors.Wrap(err, "executing template")
}
