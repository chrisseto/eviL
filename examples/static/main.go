package main

import (
	_ "embed"
	"fmt"
	"log"
	"net"
	"net/http"
	stdTemplate "text/template"

	"github.com/chrisseto/evil"
	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
)

//go:embed index.html
var indexHTML string

type Static struct{}

func (Static) OnMount(*channel.Session) error {
	return nil
}

func (Static) HandleEvent(*channel.Session, *channel.Event) error {
	return nil
}

func (Static) Template() *template.Template {
	return template.Compile(stdTemplate.Must(
		stdTemplate.New("Static").Parse(
			`Hello, World!`,
		),
	))
}

func main() {
	srv := evil.NewServer([]byte(`some secret`))
	srv.Mount("/", Static{})

	lis, err := net.Listen("tcp", "127.0.0.1:5656")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s...\n", lis.Addr())

	mux := http.NewServeMux()
	mux.Handle("/live/websocket", srv)
	mux.Handle("/static/", http.StripPrefix("/static/", evil.StaticHandler))
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		content, err := srv.RenderView("Static")
		if err != nil {
			panic(err)
		}

		rw.Write([]byte(fmt.Sprintf(indexHTML, content)))
	})

	if err := http.Serve(lis, mux); err != nil {
		log.Fatal(err)
	}
}
