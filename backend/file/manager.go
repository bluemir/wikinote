package file

import (
	"io/ioutil"
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
	return ioutil.WriteFile(m.GetFullPath(path), data, 0644)
}
func (m *manager) GetFullPath(path string) string {
	return filepath.Join(m.basepath, path)
}
