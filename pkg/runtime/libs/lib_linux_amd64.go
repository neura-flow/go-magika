//go:build linux && amd64

package libs

import (
	"embed"
	_ "embed"
)

//go:embed all:linux-amd64
var libFS embed.FS
