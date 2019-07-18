package backend

import (
	"html/template"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/pkgs/auth"
)

type PluginClause interface {
	Footer(path string) []PluginResult
	PreSave(path string, data []byte) ([]byte, error)
	PostSave(path string, data []byte) error
	OnReadWiki(ctx *auth.Context, path string, data []byte) ([]byte, error)
	AuthCheck(ctx *auth.Context) (auth.Result, error)
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
func (b *pluginClause) PreSave(path string, data []byte) ([]byte, error) {
	attr := b.File().Attr(path)
	d := data
	for _, p := range b.plugins.preSave {
		var err error
		d, err = p.OnPreSave(path, d, attr)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
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
func (b *pluginClause) OnReadWiki(ctx *auth.Context, path string, data []byte) ([]byte, error) {
	//attr := b.File().Attr(path)
	d := data
	for _, p := range b.plugins.onReadWiki {
		var err error
		d, err = p.OnReadWiki(ctx, path, d)
		if err != nil {
			return data, err
		}
	}
	return d, nil
}

func (b *pluginClause) AuthCheck(ctx *auth.Context) (auth.Result, error) {
	for _, plugin := range b.plugins.authz {
		result, err := plugin.AuthCheck(ctx)
		if err != nil {
			return auth.Unknown, err
		}
		if result == auth.Reject {
			return auth.Reject, nil
		}
	}
	return auth.Accept, nil
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
