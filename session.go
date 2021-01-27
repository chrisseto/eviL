package evil

import (
	"encoding/base64"
	"encoding/json"
	"sync"

	"github.com/cockroachdb/errors"
)

// Session is an evil session. NOT a website session
type Session struct {
	ID   string   `json:"id"`
	View string   `json:"view"`
	data sync.Map `json:"-"`
}

func NewSession(view string) Session {
	return Session{
		ID:   ID(),
		View: view,
		data: sync.Map{},
	}
}

func (s *Session) Get(key interface{}) (interface{}, bool) {
	return s.data.Load(key)
}

func (s *Session) Set(key, value interface{}) {
	s.data.Store(key, value)
}

func (s *Session) Encode() string {
	out, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(out)
}

func DecodeSession(data string) (*Session, error) {
	out, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(out, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

type SessionFactory struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

func NewSessionFactory() *SessionFactory {
	return &SessionFactory{
		sessions: map[string]*Session{},
	}
}

func (f *SessionFactory) NewSession(view string) (*Session, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	s := NewSession(view)
	f.sessions[s.ID] = &s

	return &s, nil
}

func (f *SessionFactory) LoadSession(id string) (*Session, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if s, ok := f.sessions[id]; ok {
		return s, nil
	}

	return nil, errors.Newf("no such session: %s", id)
}
