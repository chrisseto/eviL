package evil

import (
	// "html/template"
	// "net"
	"net/http"

	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
	"github.com/dgrijalva/jwt-go"
)

type Server struct {
	hub     *channel.Hub
	channel *LiveViewChannel
	secret  []byte
}

func NewServer(
	tpl *template.Template,
	secret []byte,
) *Server {
	hub := channel.NewHub()

	channel := LiveViewChannel{
		SessionFactory: NewSessionFactory(),
		Template:       tpl,
		Views:          map[string]View{},
	}

	hub.Register("lv:*", &channel)

	return &Server{
		hub:     hub,
		secret:  secret,
		channel: &channel,
	}
}

func (s *Server) NewToken(view string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &SessionClaims{
		View: view,
	})
	return token.SignedString(s.secret)
}

func (s *Server) RegisterView(name string, view View) {
	s.channel.RegisterView(name, view)
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.hub.ServeHTTP(rw, r)
}
