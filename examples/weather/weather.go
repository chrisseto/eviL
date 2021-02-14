package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/chrisseto/evil"
	"github.com/chrisseto/evil/channel"
)

type WeatherView struct{}

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
	weather := s.Get("weather")
	location := s.Get("location")

	return map[string]interface{}{
		"Location": location,
		"Weather":  weather,
	}, nil
}

func (v *WeatherView) HandleEvent(s *evil.Session, e *channel.Event) error {
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
