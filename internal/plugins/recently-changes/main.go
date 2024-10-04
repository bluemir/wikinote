package recent

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/plugins"
	"github.com/bluemir/wikinote/internal/pubsub"
)

type Options struct {
	Limit int
}
type Recents struct {
	opt   *Options
	store metadata.IStore
}
type Data struct {
	Path string
	Time time.Time
}

var defaultConfig = `
limit: 10
`

func init() {
	plugins.Register("recently-changes", New, defaultConfig, &Options{Limit: 10})
}
func New(ctx context.Context, conf any, store metadata.IStore, hub *pubsub.Hub) (plugins.Plugin, error) {
	opt, ok := conf.(*Options)
	if !ok {
		return nil, errors.Errorf("option type not matched: %T", conf)
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
	ctx := context.Background()
	str, err := r.store.Take(ctx, ".plugins", "recently-changed")
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
	ctx := context.Background()
	buf, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return r.store.Save(ctx, ".plugins", "recently-changed", string(buf))
}

func (*Recents) SetConfig(ctx context.Context, conf any) error {
	return nil
}
