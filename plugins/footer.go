package plugins

import (
	"html/template"
)

type FooterPlugin interface {
	Footer(path string) (template.HTML, error)
}
