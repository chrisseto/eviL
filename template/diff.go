package template

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Diff struct {
	Static     []string
	Dynamic    []interface{}
	Components []string
	// Title *string
}

var _ fmt.Stringer = &Diff{}
var _ json.Marshaler = &Diff{}

func (d Diff) String() string {
	var buf strings.Builder

	for i, s := range d.Static {
		buf.WriteString(s)
		if i < len(d.Dynamic) {
			switch s := d.Dynamic[i].(type) {
			case string:
				buf.WriteString(s)
			case fmt.Stringer:
				buf.WriteString(s.String())
			}
		}
	}

	return buf.String()
}

func (d Diff) MarshalJSON() ([]byte, error) {
	out := map[string]interface{}{}

	// TODO there's a guarantee about static length vs
	// dynamic length in phoenix somewhere.
	out["s"] = d.Static

	for i := range d.Dynamic {
		out[strconv.FormatInt(int64(i), 10)] = d.Dynamic[i]
	}

	return json.Marshal(out)
}
