package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Manager interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
	Search(query string) (interface{}, error)

	GetFullPath(path string) string
}

func New(basepath string) (Manager, error) {
	return &manager{
		basepath: basepath,
	}, nil
}

type manager struct {
	basepath string
}

func (m *manager) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(m.GetFullPath(path))
}
func (m *manager) Write(path string, data []byte) error {
	fullpath := m.GetFullPath(path)
	dirpath := filepath.Dir(fullpath)
	if err := os.MkdirAll(dirpath, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(m.GetFullPath(path), data, 0644)
}
func (m *manager) GetFullPath(path string) string {
	return filepath.Join(m.basepath, path)
}
