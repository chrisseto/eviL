package evil

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
	"github.com/cockroachdb/datadriven"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

type Weather struct{}

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
	tpl, err := template.ParseGlob("testdata/thermostat.d/*")
	require.NoError(t, err)

	srv := NewServer(tpl, []byte("password123"))

	srv.RegisterView("thermostat.gohtml", &Thermostat{})

	s := httptest.NewServer(srv)
	defer s.Close()

	conn, _, err := websocket.DefaultDialer.Dial(
		strings.Replace(s.URL, "http", "ws", 1),
		nil,
	)
	require.NoError(t, err)

	defer conn.Close()

	token, err := srv.NewToken("thermostat.gohtml")
	require.NoError(t, err)

	datadriven.RunTest(t, "testdata/thermostat", func(t *testing.T, d *datadriven.TestData) string {
		switch d.Cmd {
		case "send":
			message := strings.ReplaceAll(d.Input, "$TOKEN", token)

			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
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
