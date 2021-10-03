package dist

import "embed"

//go:embed *
var StaticFiles embed.FS
