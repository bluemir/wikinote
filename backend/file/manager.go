package file

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Manager interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
	List(path string) ([]os.FileInfo, error)
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
func (m *manager) List(path string) ([]os.FileInfo, error) {
	p := m.GetFullPath(path)
	ext := filepath.Ext(path)
	p = p[:len(p)-len(ext)]
	logrus.Debug(p)
	return ioutil.ReadDir(p)
}
