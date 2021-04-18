package evil

import (
	"github.com/chrisseto/evil/channel"
)

type View interface {
	OnMount(*channel.Session) error
	HandleEvent(*channel.Session, *channel.Event) error
}
