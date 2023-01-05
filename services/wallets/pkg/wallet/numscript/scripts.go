//nolint:golint
package numscript

import (
	"bytes"
	_ "embed"
	"text/template"
)

var (
	//go:embed confirm-hold.num
	ConfirmHold string
	//go:embed cancel-hold.num
	CancelHold string
)

func BuildConfirmHoldScript(final bool, asset string) string {
	buf := bytes.NewBufferString("")
	tpl := template.Must(template.New("confirm-hold").Parse(ConfirmHold))
	if err := tpl.Execute(buf, map[string]any{
		"Final": final,
		"Asset": asset,
	}); err != nil {
		panic(err)
	}
	return buf.String()
}
