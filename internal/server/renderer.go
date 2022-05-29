package server

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/static"
)

func NewRenderer() (*template.Template, error) {
	tmpl := template.New("__root__")
	tmpl.Funcs(template.FuncMap{
		"join": strings.Join,
	})

	static.HTMLTemplates.Walk("/", func(path string, info os.FileInfo, err error) error {
		logrus.Tracef("find %s", path)
		if info.IsDir() && info.Name()[0] == '.' && path != "/" {
			return filepath.SkipDir
		}
		if info.IsDir() || info.Name()[0] == '.' || !strings.HasSuffix(path, ".html") {
			return nil
		}
		logrus.Debugf("parse template: path: %s", path)

		tmpl, err = tmpl.Parse(static.HTMLTemplates.MustString(path))
		if err != nil {
			return err
		}
		return nil
	})

	return tmpl, nil
}
