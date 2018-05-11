package plugins

import (
	"html/template"
)

type FooterPlugin interface {
	Footer(path string) (template.HTML, error)
}
type AfterWikiSavePlugin interface {
	AfterWikiSave(path string, data []byte) error
}
