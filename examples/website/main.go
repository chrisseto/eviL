package main

import (
	_ "embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
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

func httpRedirector(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") == "http" {
			r.URL.Scheme = "https"
			r.URL.Host = "evilview.dev"
			http.Redirect(rw, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:5757"
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = `sup3r5ecr3t`
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	srv := evil.NewServer([]byte(secret))
	srv.Mount("/", Clock{})

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

	var handler http.Handler = mux
	if os.Getenv("REDIRECT_HTTPS") != "" {
		handler = httpRedirector(handler)
	}

	log.Printf("listening on %s...\n", lis.Addr())
	if err := http.Serve(lis, handler); err != nil {
		log.Fatal(err)
	}
}
