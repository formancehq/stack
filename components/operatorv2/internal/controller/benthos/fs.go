package benthos

import (
	"embed"
)

//go:embed global
var Global embed.FS

//go:embed audit
var Audit embed.FS
