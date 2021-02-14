package main

import (
	"context"
	"time"

	"github.com/chrisseto/evil"
	"github.com/chrisseto/evil/channel"
)

func Every(ctx context.Context, tick time.Duration, fn func() error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(tick):
			if err := fn(); err != nil {
				return err
			}
		}
	}
}

type ClockView struct{}

func (v *ClockView) OnMount(s *evil.Session) error {
	s.Set("Time", time.Now().Format(time.RFC1123))
	go Every(s.Context(), time.Second, func() error {
		s.Set("Time", time.Now().Format(time.RFC1123))
		return nil
	})
	return nil
}

func (v *ClockView) ToArgs(s *evil.Session) (interface{}, error) {
	return map[string]interface{}{
		"Time": s.Get("Time"),
	}, nil
}

func (v *ClockView) HandleEvent(s *evil.Session, e *channel.Event) error {
	return nil
}
