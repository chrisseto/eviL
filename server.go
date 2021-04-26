package evil

import (
	"net/http"

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
		Views:    make(map[string]View),
		Sessions: make(map[string]*session),
		// TODO find a better way to pass this in?
		broadcast: hub.Broadcast,
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
	s.channel.RegisterView(view)
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.hub.ServeHTTP(rw, r)
}

func (s *Server) RenderView(viewName string) (string, error) {
	view := s.channel.Views[viewName]

	instance, err := s.channel.SpawnInstance(ID(), view)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, instance.Claims())
	sessionToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}

	diff, err := view.Template().Execute(instance.assigns)
	if err != nil {
		return "", err
	}

	// TODO allow tag to be customizable
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
