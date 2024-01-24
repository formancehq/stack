package ledgers

import (
	"embed"
	_ "embed"
)

//go:embed assets/Caddyfile.gotpl
var Caddyfile string

//go:embed assets/reindex
var reindexStreams embed.FS
