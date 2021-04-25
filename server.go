package evil

import (
	// "net"
	"net/http"
	"reflect"

	"github.com/chrisseto/evil/channel"
	// "github.com/chrisseto/evil/template"

	// "github.com/cockroachdb/errors"
	"github.com/dgrijalva/jwt-go"
)

type Server struct {
	hub     *channel.Hub
	channel *LiveViewChannel
	secret  []byte
}

func NewServer(
	secret []byte,
) *Server {
	hub := channel.NewHub()

	channel := LiveViewChannel{
		Secret:   secret,
		Views:    map[string]View{},
	}

	hub.Register("lv:*", &channel)

	return &Server{
		hub:     hub,
		secret:  secret,
		channel: &channel,
	}
}

func (s *Server) NewToken(id, view string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &SessionClaims{
		ID:   id,
		View: view,
	})
	return token.SignedString(s.secret)
}

func (s *Server) Mount(path string, view View) {
	// TODO use path
	name := reflect.TypeOf(view).Name()
	s.channel.RegisterView(name, view)
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.hub.ServeHTTP(rw, r)
}

func (s *Server) RenderView(viewName string) (string, error) {
	instance, err := s.channel.SpawnInstance(viewName)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, instance.Claims())
	sessionToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}

	diff, err := s.channel.Views[viewName].Template().Execute(nil)
	if err != nil {
		return "", err
	}

	return RenderTag(
		"div",
		map[string]string{
			"id":               instance.ID,
			"data-phx-main":    "true",
			"data-phx-static":  "TODO",
			"data-phx-session": sessionToken,
			"data-phx-view":    viewName,
		},
		diff.String(),
	)
}
