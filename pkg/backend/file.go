package backend

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bluemir/wikinote/pkg/fileattr"
)

func (backend *Backend) Object(path string) (map[string]string, error) {
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

	return result, nil
}

func (backend *Backend) FileRead(path string) ([]byte, error) {
	// data, err := backend.plugins.triggerFileReadHook(path, data)
	return ioutil.ReadFile(backend.GetFullPath(path))
}
func (backend *Backend) FileWrite(path string, data []byte) error {
	data, err := backend.Plugin.TriggerFileWriteHook(path, data)
	if err != nil {
		return err
	}

	fullpath := backend.GetFullPath(path)
	dirpath := filepath.Dir(fullpath)
	if err := os.MkdirAll(dirpath, 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(backend.GetFullPath(path), data, 0644); err != nil {
		return err
	}

	return nil
}

func (backend *Backend) GetFullPath(path string) string {
	return filepath.Join(backend.Config.Wikipath, path)
}
