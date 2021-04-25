package evil

import (
	"testing"

	"github.com/cockroachdb/datadriven"
)

func TestRenderTag(t *testing.T) {
	datadriven.RunTest(t, "testdata/render_tag", func(t *testing.T, d *datadriven.TestData) string {
		switch d.Cmd {
		case "render":
			tag := d.CmdArgs[0].Key
			attrs := make(map[string]string, len(d.CmdArgs)-1)
			for _, cmd := range d.CmdArgs[1:] {
				attrs[cmd.Key] = cmd.Vals[0]
			}
			data, err := RenderTag(tag, attrs, d.Input)
			if err != nil {
				return err.Error()
			}
			return data

		default:
			d.Fatalf(t, "unknown command: %s", d.Cmd)
			return d.Expected
		}
	})
}
