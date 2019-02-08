package backend

import (
	"bytes"

	//"gopkg.in/russross/blackfriday.v2"
	//"github.com/russross/blackfriday"
	"github.com/russross/blackfriday/v2"
)

func (b *backend) Render(input []byte) ([]byte, error) {
	return blackfriday.Run(
		bytes.Replace(input, []byte("\r\n"), []byte("\n"), -1),
		blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.Footnotes|blackfriday.AutoHeadingIDs),
		blackfriday.WithRenderer(
			blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
				Flags: blackfriday.CommonHTMLFlags | blackfriday.HrefTargetBlank, //| TOC
			}),
		),
	), nil
}
