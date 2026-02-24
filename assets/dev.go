//go:build !prod

package assets

import "os"

func init() {
	Static = &mappedFS{fs: os.DirFS("assets"), prefix: src}
	HTMLTemplates = os.DirFS("assets/html-templates")
}
