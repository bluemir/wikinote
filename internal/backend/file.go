package backend

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bluemir/wikinote/internal/fileattr"
)

func (backend *Backend) Object(path string) (map[string]string, error) {
	if strings.HasPrefix(path, "/!/") {
		return map[string]string{"kind": "special"}, nil
	}
	// TODO make kind from ext
	attrs, err := backend.FileAttr.Find(&fileattr.FileAttr{
		Path: path,
	})

	if err != nil {
		return map[string]string{}, err
	}

	result := map[string]string{}
	for _, attr := range attrs {
		result[attr.Key] = attr.Value
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".md":
		result["kind"] = "wiki"
	}

	return result, nil
}

func (backend *Backend) FileRead(path string) ([]byte, error) {
	// data, err := backend.plugins.triggerFileReadHook(path, data)
	return ioutil.ReadFile(backend.getFullPath(path))
}

func (backend *Backend) FileWrite(path string, data []byte) error {
	data, err := backend.Plugin.TriggerFileWriteHook(path, data)
	if err != nil {
		return err
	}

	fullpath := backend.getFullPath(path)
	dirpath := filepath.Dir(fullpath)
	if err := os.MkdirAll(dirpath, 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(backend.getFullPath(path), data, 0644); err != nil {
		return err
	}

	return nil
}

func (backend *Backend) FileDelete(path string) error {
	return os.Remove(backend.getFullPath(path))
}

func (backend *Backend) getFullPath(path string) string {
	return filepath.Join(backend.Config.Wikipath, path)
}
