package channel

import (
	"encoding/json"

	"github.com/olahol/melody"
)

func writeJSON(session *melody.Session, data interface{}) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return session.Write(out)
}
