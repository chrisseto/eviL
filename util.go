package evil

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"text/template"
)

func ID() string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	var id strings.Builder

	id.WriteString("phx-")

	for i := 0; i < 5; i++ {
		id.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}

	return id.String()
}

func htmlAttributes(attrs map[string]string) string {
	keys := make([]string, 0, len(attrs))
	for key := range attrs {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var buf strings.Builder
	for i, key := range keys {
		if i > 0 {
			buf.WriteRune(' ')
		}

		fmt.Fprintf(&buf, `%s="%s"`, key, attrs[key])
	}

	return buf.String()
}

var tpl = template.Must(template.
	New("").
	Funcs(template.FuncMap{
		"attrs": htmlAttributes,
	}).
	Parse(`<{{ .Tag }} {{ attrs .Attrs }}>
	{{ .Content }}
</{{ .Tag }}>
`))

func RenderTag(tag string, attrs map[string]string, content string) (string, error) {
	// TODO use net/html to construct this
	var buf strings.Builder
	if err := tpl.Execute(&buf, map[string]interface{}{
		"Tag":     tag,
		"Attrs":   attrs,
		"Content": content,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
