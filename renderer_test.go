package evil

import (
	"io/ioutil"
	"strings"
	"testing"
	"text/template"
	"text/template/parse"

	"github.com/stretchr/testify/require"
)

func LoadTestData(name string) string {
	data, err := ioutil.ReadFile("testdata/" + name)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type diff struct {
	Static     []string
	Dynamic    []string
	Components []diff
}

func (d *diff) String() string {
	var buf strings.Builder

	longest := len(d.Static)
	if len(d.Dynamic) > longest {
		longest = len(d.Dynamic)
	}

	for i := 0; i < longest; i++ {
		if len(d.Static) > i {
			buf.WriteString(d.Static[i])
		}
		if len(d.Dynamic) > i {
			buf.WriteString(d.Dynamic[i])
		}
	}

	return buf.String()
}

func ExecNode(node parse.Node, data interface{}) string {
	t := parse.New("")
	t.Root = &parse.ListNode{
		Nodes: []parse.Node{node},
	}
	tpl := template.New("")
	tpl.Tree = t
	var buf strings.Builder
	if err := tpl.Execute(&buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

func Render(root *parse.ListNode, data interface{}) *diff {
	d := diff{}
	for _, node := range root.Nodes {
		switch n := node.(type) {
		case *parse.TextNode:
			d.Static = append(d.Static, ExecNode(n, data))
		default:
			d.Dynamic = append(d.Dynamic, ExecNode(n, data))
		}
	}
	return &d
}

func TestRender(t *testing.T) {
	testCases := []struct {
		Tpl  string
		Data interface{}
	}{
		{Tpl: "<h1>{{ . }}</h1>", Data: nil},
		{Tpl: "<h1>{{ . }}</h1>", Data: "Hello World"},
		{Tpl: "<h1>{{ .Greeting }}, {{ .Subject }}</h1>", Data: map[string]string{
			"Greeting": "Hello",
			"Subject":  "World",
		}},
		{Tpl: "{{range .}}{{.}}{{end}}", Data: []int{1, 2, 3}},
		{Tpl: LoadTestData("twitter.gohtml"), Data: []int{1, 2, 3}},
	}

	for _, tc := range testCases {
		tpl, err := template.New("").Parse(tc.Tpl)
		require.NoError(t, err)

		tree, err := parse.Parse("", tc.Tpl, "{{", "}}")
		require.NoError(t, err)

		diff := Render(tree[""].Root, tc.Data)

		var buf strings.Builder
		require.NoError(t, tpl.Execute(&buf, tc.Data))

		require.Equal(t, buf.String(), diff.String())
	}
}
