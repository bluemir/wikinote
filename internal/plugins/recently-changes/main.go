package recent

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/bluemir/wikinote/internal/plugins"
)

const (
	KeyLastModified = "wikinote.bluemir.me/last-modified"
)

type Options struct {
	Limit int
}
type Recents struct {
	opt   *Options
	store *plugins.Store
}

func init() {
	plugins.Register("recent-changes", New, &Options{Limit: 10})
}
func New(o interface{}, store *plugins.Store) (plugins.Plugin, error) {
	opt, ok := o.(*Options)
	if !ok {
		return nil, errors.Errorf("option not matched")
	}

	return &Recents{opt, store}, nil
}

func (r *Recents) FileWriteHook(path string, data []byte) ([]byte, error) {
	// register lastModified
	if err := r.store.Save(&plugins.FileAttr{
		Path:  path,
		Key:   KeyLastModified,
		Value: time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		return data, err
	}

	return data, nil
}
func (r *Recents) Footer(path string) ([]byte, error) {
	attr, err := r.store.Take(&plugins.FileAttr{
		Path: path,
		Key:  KeyLastModified,
	})
	if err != nil {
		if plugins.IsNotFound(err) {
			return []byte(""), nil
		}
		return []byte{}, err
	}
	t, err := time.Parse(time.RFC3339, attr.Value)
	if err != nil {
		return []byte{}, err
	}

	return []byte("last update: " + t.Local().Format(time.RFC3339)), nil
}
func (r *Recents) Route(app gin.IRouter) error {
	app.GET("/", func(c *gin.Context) {
		attrs, err := r.store.Search(&plugins.FileAttr{
			Key: KeyLastModified,
		}, &plugins.ListOption{
			Order: "value desc",
			Limit: 10,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		result := ""
		for _, attr := range attrs {

			t, err := time.Parse(time.RFC3339, attr.Value)
			if err != nil {
				t = time.Now()
			}
			result += fmt.Sprintf(`<p><a href="%s">%s</a> %s</p>`, attr.Path, attr.Path, t.Local().Format(time.RFC3339))
		}

		c.HTML(http.StatusOK, "/view/markdown.html", gin.H{
			"data": template.HTML(result),
		})
	})

	return nil
}
