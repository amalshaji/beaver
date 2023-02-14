package web

import "embed"

//go:embed dist/index.html
var DistIndex []byte

//go:embed dist/assets
var DistAssets embed.FS
