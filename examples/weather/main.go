package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	stdTemplate "text/template"
	"time"

	"github.com/chrisseto/evil"
	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
)

//go:embed index.html
var indexHTML string

//go:embed weather.html
var weatherHTML string

type Weather struct{}

func (v Weather) OnMount(s evil.Session) error {
	s.Set("Location", "Austin")
	s.Set("Weather", "...")

	go func() {
		time.Sleep(time.Second)

		weather, err := v.getWeather("Austin")
		if err != nil {
			panic(err)
		}
		s.Set("Location", "Austin")
		s.Set("Weather", weather)
	}()

	return nil
}

func (v Weather) HandleEvent(s evil.Session, e *channel.Event) error {
	values, err := url.ParseQuery(e.Value)
	if err != nil {
		return err
	}
	weather, err := v.getWeather(values.Get("location"))
	if err != nil {
		return err
	}
	s.Set("Location", values.Get("location"))
	s.Set("Weather", weather)
	return nil
}

func (Weather) Template() *template.Template {
	return template.Compile(stdTemplate.Must(
		stdTemplate.New("").Parse(weatherHTML),
	))
}

func (v Weather) getWeather(location string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(
		"http://wttr.in/%s?format=1",
		url.PathEscape(location),
	))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func main() {
	srv := evil.NewServer([]byte(`sup3r5ecr3t!1`))
	srv.Mount("/", Weather{})

	lis, err := net.Listen("tcp", "127.0.0.1:5858")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s...\n", lis.Addr())

	mux := http.NewServeMux()
	mux.Handle("/live/websocket", srv)
	mux.Handle("/static/", http.StripPrefix("/static/", evil.StaticHandler))
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		content, err := srv.RenderView("Weather")
		if err != nil {
			panic(err)
		}

		rw.Write([]byte(fmt.Sprintf(indexHTML, content)))
	})

	if err := http.Serve(lis, mux); err != nil {
		log.Fatal(err)
	}
}
