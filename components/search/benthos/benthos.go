package benthos

import (
	"embed"
)

//go:embed resources
var Resources embed.FS

//go:embed templates
var Templates embed.FS

//go:embed streams
var Streams embed.FS
