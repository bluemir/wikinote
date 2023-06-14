//ago:build embed

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
