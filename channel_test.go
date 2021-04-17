package evil

import (
	"testing"

	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
	"github.com/cockroachdb/datadriven"
)

func ddServer(glob string) func(t *testing.T, d *datadriven.TestData) string {
	hub := channel.NewHub()

	tpl := template.Must(template.ParseGlob(glob))

	hub.Register("lv:*", &LiveViewChannel{
		Templates: nil
		SessionFactory: sf,
	})

	return func(t *testing.T, d *datadriven.TestData) string {
		switch d.Cmd {
		case "send":
			// Do stuff
		}

		return d.Expected
	}
}
