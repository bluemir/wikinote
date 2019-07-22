package draft

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v3"

	"github.com/bluemir/wikinote/plugins"
)

func init() {
	plugins.Register("draft", New)
}

func New(core plugins.Core, confBuf []byte) (plugins.Plugin, error) {
	opts := &Options{}

	if err := yaml.Unmarshal(confBuf, opts); err != nil {
		return nil, err
	}
	return &Draft{core, opts}, nil
}

type Options struct {
	User struct {
		Key   string `yaml:"key"`
		Value string `yaml:"value"`
	} `yaml:"user"`
}
type Draft struct {
	plugins.Core
	options *Options
}

func (draft *Draft) OnPreSave(ctx *plugins.AuthContext, path string, buf []byte) ([]byte, error) {
	if ctx.Subject.Attr(draft.options.User.Key) == draft.options.User.Value {
		origin, err := draft.File().Read(path)
		if err != nil {
			return buf, err
		}

		if err := draft.File().Write(draft.path(path), buf); err != nil {
			return buf, err
		}

		return origin, nil
	}
	return buf, nil
}
func (draft *Draft) OnReadWiki(ctx *plugins.AuthContext, path string, buf []byte) ([]byte, error) {
	if ctx.Subject.Attr(draft.options.User.Key) == draft.options.User.Value && ctx.Action == "edit" {
		d, err := draft.File().Read(draft.path(path))
		if err != nil {
			return buf, nil
		}

		return d, nil
	}
	return buf, nil
}

func (draft *Draft) RegisterAction(qr plugins.QueryRouter, Authz plugins.AuthzFunc) {
	qr.Register(http.MethodGet, "publish", Authz("publish"), func(c *gin.Context) {
		path := c.Request.URL.Path
		buf, err := draft.File().Read(draft.path(path))
		if err != nil {
			c.HTML(http.StatusInternalServerError, "/errors/internal-error.html", map[string]string{})
			c.Abort()
			return
		}

		if err := draft.File().Write(path, buf); err != nil {
			c.HTML(http.StatusInternalServerError, "/errors/internal-error.html", map[string]string{})
			c.Abort()
			return
		}

		c.Redirect(http.StatusSeeOther, path)
	})
}
func (draft *Draft) path(path string) string {
	dir, file := filepath.Split(path)
	return fmt.Sprintf("%s/.%s.draft", dir, file)
}
