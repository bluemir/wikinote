package backend

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileAttr struct {
	Path  string
	Key   string
	Value string
}

func (backend *Backend) Object(path string) (map[string]string, error) {
	attrs := []FileAttr{}

	if err := backend.db.Where(&FileAttr{
		Path: path,
	}).Find(&attrs).Error; err != nil {
		return map[string]string{}, err
	}

	result := map[string]string{}
	for _, attr := range attrs {
		result[attr.Key] = attr.Value
	}

	return result, nil
}

func (backend *Backend) FileRead(path string) ([]byte, error) {
	return ioutil.ReadFile(backend.GetFullPath(path))
}
func (backend *Backend) FileWrite(path string, data []byte) error {
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
