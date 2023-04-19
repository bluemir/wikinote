package server

import (
	"encoding/json"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin/render"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/static"
)

func NewRender() (render.HTMLRender, error) {
	fMap := template.FuncMap{
		"base": path.Base,
		"join": strings.Join,
		"json": json.Marshal,
		"toString": func(buf []byte) string {
			return string(buf)
		},
	}
	renderer := &HTMLRender{
		templates: map[string]*template.Template{},
	}
	// layout
	layout, err := template.New("__root__").Funcs(fMap).Parse(static.HTMLTemplates.MustString("layout.html"))
	if err != nil {
		return nil, err
	}

	if err := static.HTMLTemplates.Walk("/", func(path string, info os.FileInfo, err error) error {
		logrus.Tracef("find %s", path)
		if info.IsDir() && info.Name()[0] == '.' && path != "/" {
			return filepath.SkipDir
		}

		switch { // skip condition
		case info.IsDir():
			return nil
		case info.Name()[0] == '.':
			return nil
		case !strings.HasSuffix(path, ".html"):
			return nil
		case info.Name() == "layout.html":
			return nil
		}

		logrus.Debugf("parse template: path: %s", path)

		layout, err := layout.Clone()
		if err != nil {
			return err
		}
		tmpl, err := layout.Parse(static.HTMLTemplates.MustString(path))
		if err != nil {
			return err
		}

		renderer.templates[path] = tmpl

		return nil
	}); err != nil {
		return nil, err
	}
	return renderer, nil
}

type HTMLRender struct {
	templates map[string]*template.Template
}

func (r *HTMLRender) Instance(name string, data interface{}) render.Render {

	return render.HTML{
		Template: r.templates[name],
		Name:     "__root__",
		Data:     data,
	}
}
