//go:build prod

//go:generate sh -c "cd .. && go run ./scripts/tools/import-helper"
//go:generate mkdir -p dist/js dist/css
//go:generate esbuild src/js/index.js --outdir=dist/js/ --bundle --minify --format=esm --external:lit-html --external:bm.js/bm.module.js
//go:generate esbuild src/css/page.css --outdir=dist/css/ --bundle --minify
//go:generate esbuild src/css/element.css --outdir=dist/css/ --bundle --minify

package assets

import (
	"embed"
	"io/fs"
)

//go:embed dist/js/* dist/css/* lib/* images/*
var staticFS embed.FS

//go:embed html-templates
var htmlFS embed.FS

func init() {
	Static = &mappedFS{fs: staticFS, prefix: dist}
	var err error
	HTMLTemplates, err = fs.Sub(htmlFS, "html-templates")
	if err != nil {
		panic(err)
	}
}
