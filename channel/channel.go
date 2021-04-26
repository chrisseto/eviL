package channel

import "github.com/olahol/melody"

type Channel interface {
	Join(*melody.Session, *Message) (interface{}, error)
	Handle(*melody.Session, *Message) (interface{}, error)
	Leave(*melody.Session) error
}
