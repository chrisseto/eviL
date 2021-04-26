package evil

import (
	"encoding/json"

	"github.com/chrisseto/evil/channel"
	"github.com/cockroachdb/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/olahol/melody"
)

type LiveViewChannel struct {
	Views     map[string]View
	Secret    []byte
	Sessions  map[string]*session
	broadcast func(string, string, interface{}) error
}

var _ channel.Channel = &LiveViewChannel{}

func (c *LiveViewChannel) SpawnInstance(id string, view View) (*session, error) {
	if session, ok := c.Sessions[id]; ok {
		return session, nil
	}

	session := newSession(id, view, func(event string, data interface{}) error {
		// not the best way to handle this
		return c.broadcast(id, event, data)
	})

	c.Sessions[id] = session

	return session, errors.Wrap(
		session.start(),
		"starting session",
	)
}

func (c *LiveViewChannel) RegisterView(view View) {
	// TODO error if already exists
	c.Views[viewName(view)] = view
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

func (c *LiveViewChannel) Join(session *melody.Session, m *channel.Message) (interface{}, error) {
	// TODO move join in evil
	var j channel.Join
	if err := json.Unmarshal(m.Payload, &j); err != nil {
		return nil, errors.Wrap(err, "unmarshaling into .Join")
	}

	claims, err := c.verifySession(j.Session)
	if err != nil {
		return nil, errors.Wrap(err, "verifying session")
	}

	view := c.Views[claims.View]

	instance, err := c.SpawnInstance(claims.ID, view)
	if err != nil {
		return nil, errors.Wrap(err, "spawning instance")
	}

	session.Set("id", claims.ID)
	session.Set("view", claims.View)

	diff, err := instance.RenderDiff()

	// TODO make me a struct
	return map[string]interface{}{
		"rendered": diff,
	}, errors.Wrap(err, "executing template")
}

func (c *LiveViewChannel) Handle(session *melody.Session, m *channel.Message) (interface{}, error) {
	// TODO move event in evil
	var e channel.Event
	if err := json.Unmarshal(m.Payload, &e); err != nil {
		return nil, errors.Wrap(err, "unmarshaling into an event")
	}

	instance, ok := c.Sessions[session.MustGet("id").(string)]
	if !ok {
		panic("no session")
	}

	if err := instance.RootView.HandleEvent(instance, &e); err != nil {
		return nil, errors.Wrap(err, "handling event")
	}

	diff, err := instance.RenderDiff()

	// TODO make me a struct
	// TODO exclude statics
	return map[string]interface{}{
		"diff": diff,
	}, errors.Wrap(err, "executing template")
}

func (c *LiveViewChannel) Leave(session *melody.Session) error {
	// TODO
	return nil
}
