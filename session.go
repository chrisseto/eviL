package evil

import (
	"sync"
	"time"

	"github.com/chrisseto/evil/template"
	"github.com/olahol/melody"
)

type SessionClaims struct {
	ID   string
	View string
}

func (s *SessionClaims) Valid() error {
	return nil
}

type Session interface {
	Done() <-chan struct{}
	Set(key string, value interface{})
}

type session struct {
	ID       string
	RootView View

	mu           sync.Mutex
	assigns      map[string]interface{}
	participants map[*melody.Session]bool
	started      bool
	changed      bool

	broadcast func(string, interface{}) error
	done      <-chan struct{}
}

func newSession(
	id string, rootView View, broadcast func(string, interface{}) error,
) *session {
	return &session{
		ID:        id,
		RootView:  rootView,
		broadcast: broadcast,
		done:      make(<-chan struct{}),
		assigns:   make(map[string]interface{}),
	}
}

func (s *session) Done() <-chan struct{} {
	return s.done
}

func (s *session) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// I bet a sync.condition would be great here
	s.changed = true
	s.assigns[key] = value
}

func (i *session) Claims() *SessionClaims {
	return &SessionClaims{
		ID:   i.ID,
		View: viewName(i.RootView),
	}
}

func (s *session) start() error {
	s.mu.Lock()
	if s.started {
		return nil
	}
	s.started = true
	go s.doChangeLoop()
	s.mu.Unlock()
	return s.RootView.OnMount(s)
}

func (s *session) doChangeLoop() {
	for {
		select {
		case <-s.done:
			return
		case <-time.After(time.Second):
		}

		if !s.changed {
			continue
		}

		diff, err := s.RenderDiff()
		if err != nil {
			panic(err)
		}

		if err := s.broadcast("diff", diff); err != nil {
			panic(err)
		}
	}
}

func (s *session) addParticpant(session *melody.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.participants[session] = true
}

func (s *session) removeParticpant(session *melody.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.participants, session)
}

func (s *session) RenderDiff() (template.Diff, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	diff, err := s.RootView.Template().Execute(s.assigns)
	if err != nil {
		// Should probably be broadcast this
		// Maybe just logging
		// Likely need an evil.Error marker/type/thing
		return template.Diff{}, err
	}

	s.changed = false
	return diff, nil
}
