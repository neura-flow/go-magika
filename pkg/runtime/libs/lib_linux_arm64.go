//go:build linux && arm64

package libs

import (
	"embed"
	_ "embed"
)

//go:embed all:linux-arm64
var libFS embed.FS
