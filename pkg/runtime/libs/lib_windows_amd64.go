//go:build windows && amd64

package libs

import (
	"embed"
	_ "embed"
)

//go:embed all:windows-amd64
var libFS embed.FS
