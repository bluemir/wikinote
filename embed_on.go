//go:build prod

//go:generate go run ./scripts/tools/import-helper
//go:generate mkdir -p build/static/js build/static/css
//go:generate esbuild assets/js/index.js --outdir=build/static/js/ --bundle --minify --format=esm --external:lit-html --external:bm.js/bm.module.js
//go:generate esbuild assets/css/page.css --outdir=build/static/css/ --bundle --minify
//go:generate esbuild assets/css/element.css --outdir=build/static/css/ --bundle --minify

package main

import (
	"embed"

	"github.com/bluemir/wikinote/internal/assets"
)

//go:embed build/static
//go:embed assets/html-templates
var embedFS embed.FS

func init() {
	if err := assets.InitFS(embedFS); err != nil {
		panic(err)
	}

}
