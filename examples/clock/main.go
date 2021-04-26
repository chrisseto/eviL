package main

import (
	_ "embed"
	"fmt"
	"log"
	"net"
	"net/http"
	stdTemplate "text/template"
	"time"

	"github.com/chrisseto/evil"
	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
)

//go:embed index.html
var indexHTML string

//go:embed clock.html
var clockHTML string

type Clock struct{}

func (c Clock) OnMount(s evil.Session) error {
	go c.doTick(s)
	s.Set("Time", time.Now().Format(time.RFC1123))
	return nil
}

func (Clock) doTick(s evil.Session) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.Done():
			return
		case <-ticker.C:
			s.Set("Time", time.Now().Format(time.RFC1123))
		}
	}
}

func (Clock) HandleEvent(evil.Session, *channel.Event) error {
	return nil
}

func (Clock) Template() *template.Template {
	return template.Compile(stdTemplate.Must(
		stdTemplate.New("Static").Parse(clockHTML),
	))
}

func main() {
	srv := evil.NewServer([]byte(`sup3r5ecr3t!1`))
	srv.Mount("/", Clock{})

	lis, err := net.Listen("tcp", "127.0.0.1:5757")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s...\n", lis.Addr())

	mux := http.NewServeMux()
	mux.Handle("/live/websocket", srv)
	mux.Handle("/static/", http.StripPrefix("/static/", evil.StaticHandler))
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		content, err := srv.RenderView("Clock")
		if err != nil {
			panic(err)
		}

		rw.Write([]byte(fmt.Sprintf(indexHTML, content)))
	})

	if err := http.Serve(lis, mux); err != nil {
		log.Fatal(err)
	}
}
