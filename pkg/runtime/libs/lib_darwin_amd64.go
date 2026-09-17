//go:build darwin && amd64

package libs

import (
	"embed"
	_ "embed"
)

//go:embed all:darwin-amd64
var libFS embed.FS
