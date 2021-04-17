package template

import (
	// "html/template"
	// "fmt"
	// "strings"
	// "text/template"
	// "text/template/parse"
)

// func ParseGlob(pattern string) (*Template, error) {
// 	tpl, err := template.ParseGlob(pattern)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return Compile(tpl), nil
// }

// func Parse(text string) (*Template, error) {
// 	template.Template
// 	root := template.New("").Funcs(template.FuncMap{
// 		evilComponent: func(args ...interface{}) error {
// 			panic("components not allowed here")
// 		},
// 	})
// 	out, err := root.Parse(text)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return Compile(out), nil
// }

// func Compile(t *template.Template) *Template {
// 	t = t.Funcs(template.FuncMap{
// 		evilComponent: func(name string, args ...interface{}) (string, error) {
// 			var buf strings.Builder
// 			if err := t.ExecuteTemplate(&buf, name, args); err != nil {
// 				return "", err
// 			}
// 			return buf.String(), nil
// 		},
// 	})

// 	out := &Template{
// 		root:       t,
// 		components: map[string]*Template{},
// 	}

// 	for _, node := range t.Tree.Root.Nodes {
// 		switch n := node.(type) {
// 		case *parse.TextNode:
// 			out.static = append(out.static, string(n.Text))
// 		case *parse.ActionNode:
// 			if isEvilComponent(n) {
// 				name := n.Pipe.Cmds[0].Args[1].(*parse.StringNode).Text
// 				tpl := t.Lookup(name)
// 				if tpl == nil {
// 					panic(fmt.Sprintf("no such template: %s", name))
// 				}
// 				out.components[name] = Compile(tpl)
// 				// parse.ActionNode
// 				// out.dynamic = append(out.dynamic, Component{
// 				// })
// 				// panic("You done it")
// 				out.dynamic = append(out.dynamic, n)
// 			} else {
// 				out.dynamic = append(out.dynamic, n)
// 			}
// 		default:
// 			out.dynamic = append(out.dynamic, n)
// 		}
// 	}

// 	return out
// }
