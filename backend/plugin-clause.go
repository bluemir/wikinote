package backend

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

type PluginClause interface {
	Footer(path string) []PluginResult
	AfterWikiSave(path string, data []byte) error
	RegisterRouter(r gin.IRouter)
}

type pluginClause struct {
	*backend
}

func (b *pluginClause) Footer(path string) []PluginResult {
	result := []PluginResult{}
	for _, p := range b.plugins.footer {
		d, e := p.Footer(path)
		result = append(result, PluginResult{
			d, e,
		})
	}
	return result
}
func (b *pluginClause) AfterWikiSave(path string, data []byte) error {
	for _, p := range b.plugins.afterWikiSave {
		e := p.AfterWikiSave(path, data)
		if e != nil {
			return e
		}
	}
	return nil
}
func (b *pluginClause) RegisterRouter(r gin.IRouter) {
	for name, p := range b.plugins.registerRouter {
		p.RegisterRouter(r.Group(name))
	}
}

type PluginResult struct {
	Data template.HTML
	Err  error
}
