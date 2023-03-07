package backend

import (
	"bytes"

	"github.com/russross/blackfriday/v2"
)

func (backend *Backend) Render(input []byte) ([]byte, error) {
	return blackfriday.Run(
		bytes.Replace(input, []byte("\r\n"), []byte("\n"), -1),
		blackfriday.WithExtensions(
			0|
				blackfriday.CommonExtensions|
				blackfriday.Footnotes|
				blackfriday.FencedCode|
				blackfriday.Autolink|
				blackfriday.AutoHeadingIDs|
				blackfriday.Tables,
		),
		blackfriday.WithRenderer(
			blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
				Flags: blackfriday.CommonHTMLFlags | blackfriday.HrefTargetBlank, //| TOC
			}),
		),
	), nil
}
