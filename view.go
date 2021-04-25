package evil

import (
	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
)

type View interface {
	Template() *template.Template
	OnMount(*channel.Session) error
	HandleEvent(*channel.Session, *channel.Event) error
}
