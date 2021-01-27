package evil

import (
	"html/template"
	// "github.com/cockroachdb/errors"
	// "github.com/gorilla/websocket"
)

type Server struct {
	Template *template.Template
	Views    map[string]View
}

func NewServer(views map[string]View) *Server {
	s := &Server{
		// Template: registerFuncs(template.New("")),
		Views: views,
	}

	return s
}

// func (s *Server) Run(ws *websocket.Conn) error {
// 	for {
// 		var msg Message
// 		if err := ws.ReadJSON(&msg); err != nil {
// 			return errors.Wrap(err, "unmarshaling message")
// 		}

// 		resp, err := s.handleMessage(&msg)
// 		if err != nil {
// 			return errors.Wrap(err, "handling message")
// 		}

// 		if err := ws.WriteJSON(resp); err != nil {
// 			return errors.Wrap(err, "writing response")
// 		}
// 	}
// }
