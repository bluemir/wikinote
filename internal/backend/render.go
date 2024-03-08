package backend

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func (backend *Backend) Render(input []byte) ([]byte, error) {

	extensions := 0 |
		parser.CommonExtensions |
		parser.Footnotes |
		parser.FencedCode |
		parser.Autolink |
		parser.AutoHeadingIDs |
		parser.Tables

	return markdown.Render(
		parser.NewWithExtensions(extensions).Parse(input),
		html.NewRenderer(html.RendererOptions{
			Flags: html.CommonFlags | html.HrefTargetBlank,
		}),
	), nil
}
