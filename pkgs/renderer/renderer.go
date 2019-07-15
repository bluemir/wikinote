package renderer

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin/render"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkgs/dist"
)

type LayoutRenderer struct {
	templates map[string]*template.Template
}

func NewRenderer() render.HTMLRender {
	templates := map[string]*template.Template{}
	box := dist.HTMLs
	layout := template.New("layout")
	box.Walk("/_layout", func(path string, info os.FileInfo, err error) error {
		logrus.Debugf("[parse layout] name :%s, path: %s", info.Name(), path)
		if (info.Name()[0] == '.') && path != "/" {
			return filepath.SkipDir
		}
		if info.IsDir() || path[0] == '.' || !strings.HasSuffix(path, ".html") {
			return nil
		}

		layout = template.Must(layout.Parse(box.MustString(path)))
		return nil
	})
	layout.Funcs(template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"title": strings.Title,
	})

	logrus.Debug(layout.DefinedTemplates())
	box.Walk("/", func(path string, info os.FileInfo, err error) error {
		logrus.Debugf("[parse template] name :%s, path: %s", info.Name(), path)
		if (info.Name()[0] == '.' || info.Name()[0] == '_') && path != "/" {
			return filepath.SkipDir
		}
		if info.IsDir() || path[0] == '_' || !strings.HasSuffix(path, ".html") {
			return nil
		}

		templates[path] = template.Must(template.Must(layout.Clone()).Parse(box.MustString(path)))
		return nil
	})

	keys := make([]string, 0, len(templates))
	for k := range templates {
		keys = append(keys, k)
	}
	logrus.Infof("loaded template: %+v", keys)

	return &LayoutRenderer{
		templates: templates,
	}
}

func (r *LayoutRenderer) Instance(name string, data interface{}) render.Render {
	// https://github.com/michelloworld/ez-gin-template
	// just pass all data to RenderContext
	// because I can't return error here

	return &RenderContext{
		name:     name,
		data:     data,
		template: r.templates[name],
	}
}

type RenderContext struct {
	name     string
	data     interface{}
	template *template.Template
}

func (r *RenderContext) Render(w http.ResponseWriter) error {
	data, ok := r.data.(*userData)
	if !ok {
		return fmt.Errorf("wrong type, data is not *UserData*")
	}
	if r.template == nil {
		return fmt.Errorf("template '%s' no found", r.name)
	}

	rd, err := data.Pack()
	if err != nil {
		return err
	}

	return r.template.Execute(w, rd)
}
func (r *RenderContext) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"text/html"}
	}
}
