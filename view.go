package evil

import (
	"reflect"

	"github.com/chrisseto/evil/channel"
	"github.com/chrisseto/evil/template"
)

type View interface {
	Template() *template.Template
	OnMount(Session) error
	HandleEvent(Session, *channel.Event) error
}

func viewName(view View) string {
	return reflect.TypeOf(view).Name()
}
