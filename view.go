package evil

import (
	"context"
)

type Conn struct {
	ID string
	context.Context
}

func (c *Conn) Set(key, val interface{}) {
	c.Context = context.WithValue(c.Context, key, val)
}

type ViewArgs struct {
	ID      string
	Tag     string
	View    string
	Classes []string
}

type Event struct {
	Event   string                 `json:"event"`
	Type    string                 `json:"type"`
	Uploads map[string]interface{} `json:"uploads"`
	Value   string                 `json:"value"`
}

type ViewFactory func() View

type View interface {
	OnMount(*Session) error
	ToArgs(*Session) (interface{}, error)
	HandleEvent(*Session, *Event) error
}
