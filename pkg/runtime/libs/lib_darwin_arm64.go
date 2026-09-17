//go:build darwin && arm64

package libs

import (
	"embed"
	_ "embed"
)

//go:embed all:darwin-arm64
var libFS embed.FS
