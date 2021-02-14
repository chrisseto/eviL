package main

import (
	"fmt"
	"net/http"

	"github.com/chrisseto/evil"
	"github.com/chrisseto/evil/channel"
)

func main() {
	mux := http.NewServeMux()
	hub := channel.NewHub()
	fileServer := http.FileServer(http.Dir("./dist/"))

	sf := evil.NewSessionFactory()

	renderer, err := evil.NewRenderer(
		"./dist/*.html",
		"./dist/views/*.html",
		"./dist/components/*.html",
	)
	if err != nil {
		panic(err)
	}

	renderer.SessionFactory = sf

	renderer.RegisterView("ClockView", &ClockView{})
	renderer.RegisterView("WeatherView", &WeatherView{})

	hub.Register("lv:*", &evil.LiveViewChannel{
		Renderer:       renderer,
		SessionFactory: sf,
	})

	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		if err := renderer.RenderPage(rw, "index.html"); err != nil {
			fmt.Printf("%s\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})

	mux.Handle("/live/websocket", hub)
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	if err := http.ListenAndServe("localhost:4747", mux); err != nil {
		panic(err)
	}
}
