package renderer

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin/render"
	"github.com/sirupsen/logrus"
)

type LayoutRenderer struct {
	templates map[string]*template.Template
}

func NewRenderer() render.HTMLRender {
	templates := map[string]*template.Template{}
	box := rice.MustFindBox("../../dist/html")
	layout := template.New("layout")
	box.Walk("/_layout", func(path string, info os.FileInfo, err error) error {
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
		if info.IsDir() || path[0] == '.' || path[0] == '_' || !strings.HasSuffix(path, ".html") {
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

	// Attach default data(like title)
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
	data, ok := r.data.(*UserData)
	if !ok {
		return fmt.Errorf("wrong type, data is not *UserData*")
	}
	if r.template == nil {
		return fmt.Errorf("template '%s' no found", r.name)
	}

	session := sessions.Default(data.context)

	// FIXME use const in server pkg
	username, ok := session.Get("username").(string)
	if !ok {
		username = ""
	}

	MsgInfo := session.Flashes(MSG_INFO)
	MsgWarn := session.Flashes(MSG_WARN)
	session.Save()

	return r.template.Execute(w, &RenderData{
		Path:       getWikipath(data.context),
		Breadcrumb: parseBreadcrumb(getWikipath(data.context)),
		UserName:   username,
		MsgInfo:    MsgInfo,
		MsgWarn:    MsgWarn,
		UserData:   data.data,
	})
}
func (r *RenderContext) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"text/html"}
	}
}
