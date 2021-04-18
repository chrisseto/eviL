package template

import (
	"text/template/parse"
)

const (
	evilComponent = "evil_component"
	evilRender = "evil_render"
)

func isEvilComponent(n *parse.ActionNode) bool {
	ident, ok := n.Pipe.Cmds[0].Args[0].(*parse.IdentifierNode)
	return ok && ident.Ident == evilComponent
}

func evilComponentName(n *parse.ActionNode) string {
	return n.Pipe.Cmds[0].Args[1].(*parse.StringNode).Text
}

func Must(tpl *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return tpl
}
