package template

import (
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/cockroachdb/errors"
)

// Notes:
// Template get's parsed to a tree structure
// dynamic may contain other templates
// components are just templates with special markers (effectively)
// Rendering is handled by a customized traversal function https://github.com/phoenixframework/phoenix_live_view/blob/master/lib/phoenix_live_view/diff.ex#L333
// Component state is stored in the socket

type Template struct {
	root *template.Template
}

func New() *template.Template {
	return template.New("").Funcs(template.FuncMap{
		evilComponent: func(args ...interface{}) error {
			panic("components not allowed here")
		},
		evilRender: func(name, id string) error {
			return nil
		},
	})
}

func ParseGlob(glob string) (*Template, error) {
	tpl, err := New().ParseGlob(glob)
	if err != nil {
		return nil, err
	}

	return Compile(tpl), nil
}

func Compile(stdTpl *template.Template) *Template {
	return &Template{
		root: stdTpl,
	}
}

func (t *Template) ExecuteTemplate(name string, data interface{}) (Diff, error) {
	stdTpl := t.root.Lookup(name)
	if stdTpl == nil {
		return Diff{}, errors.Newf("no such template: %s", name)
	}

	var static []string
	var dynamic []interface{}

	for _, node := range stdTpl.Tree.Root.Nodes {
		switch n := node.(type) {
		case *parse.TextNode:
			static = append(static, string(n.Text))
		case *parse.ActionNode:
			if isEvilComponent(n) {
				// Probably need some locking here?
				// name := evilComponentName(n)
				// if _, ok := t.common.components[name]; !ok {
				panic("todo")
				// }
			}
			// t.dynamic = append(t.dynamic, n)
			// 				tpl := t.Lookup(name)
			// 				if tpl == nil {
			// 					panic(fmt.Sprintf("no such template: %s", name))
			// 				}
			// 				t.common.components[name] = Compile(tpl)
			// 				// parse.ActionNode
			// 				// out.dynamic = append(out.dynamic, Component{
			// 				// })
			// 				// panic("You done it")
			// 				t.dynamic = append(t.dynamic, n)
			// 			} else {
			// 				t.dynamic = append(t.dynamic, n)
			// 			}

			data, err := t.execNode(n, data)
			if err != nil {
				return Diff{}, err
			}
			dynamic = append(dynamic, data)

		default:
			data, err := t.execNode(n, data)
			if err != nil {
				return Diff{}, err
			}
			dynamic = append(dynamic, data)
		}
	}

	return Diff{
		Static:  static,
		Dynamic: dynamic,
	}, nil
}

func (t *Template) execNode(node parse.Node, data interface{}) (interface{}, error) {
	// if n, ok := node.(*parse.ActionNode); ok && isEvilComponent(n) {
	// 	name := n.Pipe.Cmds[0].Args[1].(*parse.StringNode).Text
	// 	return t.common.components[name].Execute(data)
	// }

	// This can probably be optimized/precomputed??
	tpl, err := t.root.AddParseTree("tmpname", &parse.Tree{
		Root: &parse.ListNode{
			Nodes: []parse.Node{node},
		},
	})
	if err != nil {
		return nil, err
	}
	var buf strings.Builder
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.String(), nil
}

// type Template struct {
// 	common *common
// 	name string
// 	static     []string
// 	dynamic    []parse.Node

// 	// root       *template.Template
// 	// views      map[string]*Template
// 	// components map[string]*Template
// }

// func New(name string) *Template {
// 	tpl := &Template{name: name}
// 	tpl.init()
// 	return tpl
// }

// func (t *Template) init() {
// 	if t.common == nil {
// 		t.common = &common{
// 			root: template.New("").Funcs(template.FuncMap{
// 				evilComponent: func(args ...interface{}) error {
// 					panic("components not allowed here")
// 				},
// 			}),
// 			components: make(map[string]*Template),
// 			templates: make(map[string]*Template),
// 		}
// 	}
// }

// func (t *Template) Parse(text string) error {
// 	t.init()
// 	tpl := t.common.root.Lookup(t.name)
// 	tpl, err := tpl.Parse(text)
// 	if err != nil {
// 		return err
// 	}

// 	for _, node := range tpl.Tree.Root.Nodes {
// 		switch n := node.(type) {
// 		case *parse.TextNode:
// 			t.static = append(t.static, string(n.Text))
// 		case *parse.ActionNode:
// 			if isEvilComponent(n) {
// 				// Probably need some locking here?
// 				name := evilComponentName(n)
// 				if _, ok := t.common.components[name]; !ok {
// 					panic("todo")
// 				}
// 			}

// 			t.dynamic = append(t.dynamic, n)

// // 				tpl := t.Lookup(name)
// // 				if tpl == nil {
// // 					panic(fmt.Sprintf("no such template: %s", name))
// // 				}
// // 				t.common.components[name] = Compile(tpl)
// // 				// parse.ActionNode
// // 				// out.dynamic = append(out.dynamic, Component{
// // 				// })
// // 				// panic("You done it")
// // 				t.dynamic = append(t.dynamic, n)
// // 			} else {
// // 				t.dynamic = append(t.dynamic, n)
// // 			}
// 		default:
// 			t.dynamic = append(t.dynamic, n)
// 		}
// 	}

// 	return nil
// }

// // func (t *Template) ExecuteTemplate(name string, data interface{}) (Diff, error) {

// // }

// func (t *Template) Execute(data interface{}) (Diff, error) {
// 	dynamic := make([]interface{}, len(t.dynamic))
// 	for i, node := range t.dynamic {
// 		var err error
// 		dynamic[i], err = t.execNode(node, data)
// 		if err != nil {
// 			return Diff{}, err
// 		}
// 	}

// 	return Diff{
// 		Static:  t.static,
// 		Dynamic: dynamic,
// 	}, nil
// }

// func (t *Template) execNode(node parse.Node, data interface{}) (interface{}, error) {
// 	if n, ok := node.(*parse.ActionNode); ok && isEvilComponent(n) {
// 		name := n.Pipe.Cmds[0].Args[1].(*parse.StringNode).Text
// 		return t.common.components[name].Execute(data)
// 	}

// 	// This can probably be optimized/precomputed??
// 	tpl, err := t.common.root.AddParseTree("tmpname", &parse.Tree{
// 		Root: &parse.ListNode{
// 			Nodes: []parse.Node{node},
// 		},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	var buf strings.Builder
// 	if err := tpl.Execute(&buf, data); err != nil {
// 		return nil, err
// 	}
// 	return buf.String(), nil
// }
