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
	OnRead(path string, data []byte) ([]byte, error)
	TryRead(path string, user interface{}) error
	TryWrite(path string, user interface{}) error
	AuthCheck(ctx *auth.Context) (bool, error)
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
func (b *pluginClause) OnRead(path string, data []byte) ([]byte, error) {
	attr := b.File().Attr(path)
	d := data
	for _, p := range b.plugins.onRead {
		var err error
		d, err = p.OnRead(path, d, attr)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (b *pluginClause) TryRead(path string, user interface{}) error {
	attr := b.File().Attr(path)
	for _, plugin := range b.plugins.permission {
		err := plugin.TryRead(path, user, attr)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *pluginClause) TryWrite(path string, user interface{}) error {
	attr := b.File().Attr(path)
	for _, plugin := range b.plugins.permission {
		err := plugin.TryWrite(path, user, attr)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *pluginClause) AuthCheck(ctx *auth.Context) (bool, error) {
	for _, plugin := range b.plugins.authz {
		ok, err := plugin.AuthCheck(ctx)
		if err != nil {
			return false, err
		}

		if !ok {
			return false, nil
		}
	}
	return true, nil
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
