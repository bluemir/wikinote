package backend

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

type PluginClause interface {
	Footer(path string) []PluginResult
	PostSave(path string, data []byte) error
	RegisterRouter(r gin.IRouter)
}

type pluginClause struct {
	*backend
}

func (b *pluginClause) Footer(path string) []PluginResult {
	attr := b.File().Attr(path)
	result := []PluginResult{}
	for _, p := range b.plugins.footer {
		d, e := p.Footer(path, attr)
		result = append(result, PluginResult{
			d, e,
		})
	}
	return result
}
func (b *pluginClause) PostSave(path string, data []byte) error {
	store := b.File().Attr(path)
	for _, p := range b.plugins.postSave {
		err := p.OnPostSave(path, data, store)
		if err != nil {
			return err
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
