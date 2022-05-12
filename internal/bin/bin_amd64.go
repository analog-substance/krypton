//go:build amd64

package bin

import (
	"embed"
)

//go:embed nmap_amd64
var binFS embed.FS
