package server

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkg/dist"
)

func loadTemplates() multitemplate.Renderer {
	log := logrus.WithField("method", "server.loadTemplates")
	r := multitemplate.NewRenderer()

	layout := template.New("layout")
	dist.Templates.Walk("/_layout", func(path string, info os.FileInfo, err error) error {
		log.Debugf("[parse layout] name :%s, path: %s", info.Name(), path)

		if info.IsDir() && info.Name()[0] == '.' && path != "/" {
			return filepath.SkipDir
		}
		if info.Name()[0] == '.' {
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		layout = template.Must(layout.Parse(dist.Templates.MustString(path)))
		return nil
	})

	layout.Funcs(template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"title": strings.Title,
	})

	log.Debug(layout.DefinedTemplates())

	dist.Templates.Walk("/", func(path string, info os.FileInfo, err error) error {
		log.Debugf("[parse template] name :%s, path: %s", info.Name(), path)

		if info.IsDir() && info.Name()[0] == '.' && path != "/" {
			return filepath.SkipDir
		}
		if info.IsDir() && strings.HasPrefix(path, "/_layout") {
			return filepath.SkipDir
		}
		if info.Name()[0] == '.' {
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		log.Debugf("[register template] name :%s", path)
		r.Add(path, template.Must(template.Must(layout.Clone()).Parse(dist.Templates.MustString(path))))

		return nil
	})

	return r
}
