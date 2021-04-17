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
