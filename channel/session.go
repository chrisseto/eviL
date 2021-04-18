package channel

import (
	"context"

	"github.com/olahol/melody"
)

type Session struct {
	session *melody.Session
	ctx     context.Context
}

func (s *Session) Context() context.Context {
	return s.ctx
}

func (s *Session) Set(key string, value interface{}) {
	s.session.Set(key, value)
}

func (s *Session) Keys() map[string]interface{} {
	return s.session.Keys
}
