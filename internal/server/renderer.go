package server

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/bluemir/wikinote/internal/assets"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func NewRenderer() (*template.Template, error) {
	tmpl := template.New("__root__").Funcs(template.FuncMap{
		"base": path.Base,
		"join": func(sep string, arr []string) string {
			return strings.Join(arr, sep)
		},
		"json": json.Marshal,
		"toString": func(buf []byte) string {
			return string(buf)
		},
		"yaml": func(v any) string {
			buf, err := yaml.Marshal(v)
			if err != nil {
				logrus.Warn(err)
			}
			return string(buf)
		},
		"encode": func(v any) string {
			buf, err := json.Marshal(v)
			if err != nil {
				logrus.Warn(err)
				return ""
			}
			if bytes.Equal(buf, []byte("null")) { // it is 'null'
				return ""
			}

			return string(buf)
		},
	})

	if err := fs.WalkDir(assets.HTMLTemplates, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "read template error: path: %s", path)
		}
		logrus.Debugf("read template: path: %s", path)

		if info.IsDir() && info.Name()[0] == '.' && path != "/" {
			return filepath.SkipDir
		}
		if info.IsDir() || info.Name()[0] == '.' || !strings.HasSuffix(path, ".html") {
			return nil
		}
		logrus.Debugf("parse template: path: %s", path)

		buf, err := fs.ReadFile(assets.HTMLTemplates, path)
		if err != nil {
			return err
		}

		tmpl, err = tmpl.New(path).Parse(string(buf))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	for _, t := range tmpl.Templates() {
		logrus.Tracef("there is '%s' template", t.Name())
	}

	return tmpl, nil
}
