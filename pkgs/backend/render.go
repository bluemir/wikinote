package backend

import (
	//"github.com/russross/blackfriday"
	"gopkg.in/russross/blackfriday.v1"
	// TODO change to v2
)

func (b *backend) Render(input []byte) ([]byte, error) {
	flag := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES |
		blackfriday.HTML_HREF_TARGET_BLANK
	ext := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_FOOTNOTES |
		blackfriday.EXTENSION_AUTO_HEADER_IDS
	renderer := blackfriday.HtmlRenderer(flag, "", "")
	return blackfriday.MarkdownOptions(input, renderer, blackfriday.Options{
		Extensions: ext,
	}), nil
}
