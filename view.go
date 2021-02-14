package evil

import (
	"github.com/chrisseto/evil/channel"
)

type View interface {
	OnMount(*Session) error
	ToArgs(*Session) (interface{}, error)
	HandleEvent(*Session, *channel.Event) error
}
