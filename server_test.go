package evil

import (
	"math/rand"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
	"github.com/cockroachdb/datadriven"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

type SimpleLive struct{}

func (t *SimpleLive) OnMount(session *channel.Session) error {
	return nil
}

func (t *SimpleLive) HandleEvent(session *channel.Session, event *channel.Event) error {
	return nil
}

type Thermostat struct{}

func (t *Thermostat) OnMount(session *channel.Session) error {
	session.Set("Mode", "cooling")
	session.Set("Val", "72")
	session.Set("Time", "15:31:58")
	return nil
}

func (t *Thermostat) HandleEvent(session *channel.Session, event *channel.Event) error {
	return nil
}

func TestServer(t *testing.T) {
	rand.Seed(0)
	defer rand.Seed(time.Now().UnixNano())

	tpl, err := template.ParseGlob("testdata/simple.d/*")
	require.NoError(t, err)

	srv := NewServer(tpl, []byte("password123"))

	// srv.RegisterView("thermostat.gohtml", &Thermostat{})
	srv.RegisterView("SimpleLive", &SimpleLive{})

	s := httptest.NewServer(srv)
	defer s.Close()

	conn, _, err := websocket.DefaultDialer.Dial(
		strings.Replace(s.URL, "http", "ws", 1),
		nil,
	)
	require.NoError(t, err)

	defer conn.Close()

	var id string
	var sessionToken string

	datadriven.RunTest(t, "testdata/simple", func(t *testing.T, d *datadriven.TestData) string {
		switch d.Cmd {
		case "render":
			rendered, err := srv.RenderView(d.CmdArgs[0].Key)
			if err != nil {
				d.Fatalf(t, "%+v", err)
			}

			expr := regexp.MustCompile(`data-phx-session="([^"]+)"`)
			sessionToken = expr.FindStringSubmatch(rendered)[1]

			expr = regexp.MustCompile(`id="([^"]+)"`)
			id = expr.FindStringSubmatch(rendered)[1]

			return rendered

		case "send":
			d.Input = strings.ReplaceAll(d.Input, "$ID", id)
			d.Input = strings.ReplaceAll(d.Input, "$SESSION_TOKEN", sessionToken)

			if err := conn.WriteMessage(websocket.TextMessage, []byte(d.Input)); err != nil {
				d.Fatalf(t, "WriteMessage failed: %+v", err)
			}

		case "read":
			_, msg, err := conn.ReadMessage()
			if err != nil {
				d.Fatalf(t, "ReadMessage failed: %+v", err)
			}
			return string(msg)

		default:
			t.Fatalf("unknown command: %s", d.Cmd)
		}

		return d.Expected
	})
}
