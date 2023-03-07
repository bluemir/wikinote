package recent

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/plugins"
)

type Options struct {
	Limit int
}
type Recents struct {
	opt   *Options
	store metadata.Store
}
type Data struct {
	Path string
	Time time.Time
}

func init() {
	plugins.Register("recently-changes", New, &Options{Limit: 10})
}
func New(o interface{}, store metadata.Store) (plugins.Plugin, error) {
	opt, ok := o.(*Options)
	if !ok {
		return nil, errors.Errorf("option not matched")
	}

	return &Recents{opt, store}, nil
}

func (r *Recents) FileWriteHook(path string, data []byte) ([]byte, error) {
	list, err := r.read()
	if err != nil && !errors.Is(err, metadata.ErrNotFound) {
		return data, err
	}

	if len(list) >= r.opt.Limit {
		list = append(list[1:], Data{
			Path: path,
			Time: time.Now(),
		})
	} else {
		list = append(list, Data{
			Path: path,
			Time: time.Now(),
		})
	}

	if err := r.write(list); err != nil {
		return data, err
	}

	return data, nil
}

func (r *Recents) Route(app gin.IRouter) error {
	app.GET("/", func(c *gin.Context) {
		list, err := r.read()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		result := ""
		for _, d := range list {
			result += fmt.Sprintf(`<p><a href="%s">%s</a> %s</p>`, d.Path, d.Path, d.Time.Local().Format(time.RFC3339))
		}

		c.HTML(http.StatusOK, "/viewers/markdown.html", gin.H{
			"content": template.HTML(result),
		})
	})

	return nil
}
func (r *Recents) read() ([]Data, error) {
	str, err := r.store.Take(".plugins", "recently-changed")
	if err != nil {
		return nil, err
	}
	list := []Data{}
	if err := json.Unmarshal([]byte(str), &list); err != nil {
		return nil, err
	}
	return list, nil
}
func (r *Recents) write(list []Data) error {
	buf, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return r.store.Save(".plugins", "recently-changed", string(buf))
}
