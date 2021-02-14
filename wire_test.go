package evil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Str(s string) *string {
	return &s
}

func TestEnvelopeUnMarshal(t *testing.T) {
	testCases := []struct {
		In  string
		Out Envelope
	}{
		{
			In: `["3","3","phoenix:live_reload","phx_join",{}]`,
			Out: Envelope{
				JoinRef: Str("3"),
				Ref:     "3",
				Topic:   "phoenix:live_reload",
				Type:    TypeJoin,
				Payload: nil,
			},
		},
		{
			In: `["4","4","lv:phx-123","phx_join",{"url":"http://localhost:4000/","params":{"_csrf_token":"csrf"},"session":"<SESSION>","static":"<STATIC>","joins":0}]`,
			Out: Envelope{
				JoinRef: Str("4"),
				Ref:     "4",
				Topic:   "lv:phx-123",
				Type:    TypeJoin,
				Payload: Join{
					URL: "http://localhost:4000/",
					Params: map[string]string{
						"_csrf_token": "csrf",
					},
					Static:  "<STATIC>",
					Session: "<SESSION>",
				},
			},
		},
		// {
		// 	In: `[null,"4","phoenix","heartbeat",{}]`,
		// 	Out: Envelope{

		// 	},
		// },
	}

	for _, tc := range testCases {
		var env Envelope
		assert.NoError(t, json.Unmarshal([]byte(tc.In), &env))
		assert.Equal(t, env, tc.Out)
	}
}
