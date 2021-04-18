package template

// import (
// 	"encoding/json"
// 	"fmt"
// 	"testing"

// 	"github.com/cockroachdb/datadriven"
// 	"github.com/stretchr/testify/require"
// )

// func parseArgs(t *testing.T, d *datadriven.TestData) interface{} {
// 	if len(d.Input) == 0 {
// 		return nil
// 	}

// 	var args interface{}
// 	if err := json.Unmarshal([]byte(d.Input), &args); err != nil {
// 		d.Fatalf(t, "invalid JSON: %v", err)
// 	}

// 	return args
// }

// func ddRender() func(t *testing.T, d *datadriven.TestData) string {
// 	tpl := New("test")

// 	return func(t *testing.T, d *datadriven.TestData) string {
// 		switch d.Cmd {
// 		case "parse":
// 			if err := tpl.Parse(d.Input); err != nil {
// 				d.Fatalf(t, "failed to parse template: %v", err)
// 			}
// 			return ""
// 		case "render":
// 			diff, err := tpl.Execute(parseArgs(t, d))
// 			if err != nil {
// 				d.Fatalf(t, "failed to execute template: %v", err)
// 			}

// 			return diff.String()
// 		case "diff":
// 			diff, err := tpl.Execute(parseArgs(t, d))
// 			require.NoError(t, err)
// 			actual, err := json.MarshalIndent(&diff, "", "\t")
// 			require.NoError(t, err)
// 			return string(actual)
// 		}

// 		return d.Expected
// 	}
// }

// func TestExecute(t *testing.T) {
// 	files := []string{"static", "variables", "if-else", "components"}

// 	for _, file := range files {
// 		t.Run(file, func(t *testing.T) {
// 			datadriven.RunTest(t, fmt.Sprintf("testdata/%s", file), ddRender())
// 		})
// 	}
// }
