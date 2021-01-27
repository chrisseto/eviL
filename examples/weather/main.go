package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/chrisseto/evil"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WeatherView struct {
}

func (v *WeatherView) getWeather(location string) (string, error) {
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

func (v *WeatherView) OnMount(s *evil.Session) error {
	weather, err := v.getWeather("Austin")
	if err != nil {
		return err
	}
	s.Set("location", "Austin")
	s.Set("weather", weather)
	return nil
}

func (v *WeatherView) ToArgs(s *evil.Session) (interface{}, error) {
	weather, _ := s.Get("weather")
	location, _ := s.Get("location")

	return map[string]interface{}{
		"Location": location,
		"Weather":  weather,
	}, nil
}

func (v *WeatherView) HandleEvent(s *evil.Session, e *evil.Event) error {
	values, err := url.ParseQuery(e.Value)
	if err != nil {
		return err
	}
	weather, err := v.getWeather(values.Get("location"))
	if err != nil {
		return err
	}
	s.Set("location", values.Get("location"))
	s.Set("weather", weather)
	return nil
}

func main() {
	mux := http.NewServeMux()

	sf := evil.NewSessionFactory()

	renderer, err := evil.NewRenderer2(
		"./dist/*.html",
		"./dist/views/*.html",
		"./dist/components/*.html",
	)
	if err != nil {
		panic(err)
	}

	renderer.RegisterView("WeatherView", &WeatherView{})

	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		s, err := sf.NewSession("WeatherView")
		if err != nil {
			fmt.Printf("%s\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
		}

		if err := renderer.RenderPage(rw, "index.html", s); err != nil {
			fmt.Printf("%s\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/live/websocket", func(rw http.ResponseWriter, req *http.Request) {
		ws, err := upgrader.Upgrade(rw, req, nil)
		if err != nil {
			fmt.Printf("run: %#v\n", err)
			return
		}

		defer ws.Close()

		c := evil.Channel{
			SessionFactory: sf,
			Conn: &evil.Conn{
				Context: req.Context(),
			},
			WebSocket: ws,
			Renderer:  renderer,
		}

		if err := c.Run(); err != nil {
			fmt.Printf("run: %#v\n", err)
			return
		}
	})

	fileServer := http.FileServer(http.Dir("./dist/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	if err := http.ListenAndServe("localhost:4747", mux); err != nil {
		panic(err)
	}
}
