package evil

import (
	"encoding/json"

	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
	"github.com/cockroachdb/errors"
	"github.com/dgrijalva/jwt-go"
)

// probably going to become "session"
type Instance struct {
	ID       string
	RootView string
	assigns  map[string]interface{}
}

func (i *Instance) Get(key string) interface{} {
	return i.assigns[key]
}

func (i *Instance) Set(key string, value interface{}) {
	i.assigns[key] = value
}

func (i *Instance) Claims() *SessionClaims {
	return &SessionClaims{
		ID:   i.ID,
		View: i.RootView,
	}
}

type LiveViewChannel struct {
	Template *template.Template
	Views    map[string]View
	Secret   []byte
}

var _ channel.Channel = &LiveViewChannel{}

func (c *LiveViewChannel) SpawnInstance(rootView string) (*Instance, error) {
	if _, ok := c.Views[rootView]; !ok {
		return nil, errors.Newf("no such view: %s", rootView)
	}

	// TODO check view existance
	// TODO needs a life cycle
	return &Instance{
		ID:       ID(),
		RootView: rootView,
	}, nil
}

func (c *LiveViewChannel) RegisterView(name string, view View) {
	c.Views[name] = view
}

func (c *LiveViewChannel) verifySession(signed string) (*SessionClaims, error) {
	token, err := jwt.ParseWithClaims(signed, &SessionClaims{}, func(t *jwt.Token) (interface{}, error) {
		_ = t.Method.(*jwt.SigningMethodHMAC)
		return c.Secret, nil
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
		return nil, errors.Wrap(err, "unmarshaling into .Join")
	}

	claims, err := c.verifySession(j.Session)
	if err != nil {
		return nil, errors.Wrap(err, "verifying session")
	}

	view := c.Views[claims.View]

	// TODO this should become a "maybecallmount"
	if err := view.OnMount(session); err != nil {
		return nil, errors.Wrap(err, "calling mount")
	}

	diff, err := view.Template().Execute(session.Keys())
	if err != nil {
		return nil, errors.Wrap(err, "executing template")
	}

	// TODO make me a struct
	return map[string]interface{}{
		"rendered": diff,
	}, nil
}

func (c *LiveViewChannel) Handle(session *channel.Session, m *channel.Message) (interface{}, error) {
	var e channel.Event
	if err := json.Unmarshal(m.Payload, &e); err != nil {
		return nil, errors.Wrap(err, "unmarshaling into an event")
	}

	// if view, ok := c.Views[claims.View]; ok {
	// 	if err := view.OnMount(session); err != nil {
	// 		return nil, err
	// 	}
	// }

	diff, err := c.Template.ExecuteTemplate("", session)

	return diff, errors.Wrap(err, "executing template")
}
