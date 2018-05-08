package file

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Manager interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
	List(path string) ([]os.FileInfo, error)
	Search(query string) (interface{}, error)

	GetFullPath(path string) string
}

func New(basepath string, db *gorm.DB) (Manager, error) {
	return &manager{
		basepath: basepath,
		db:       db,
	}, nil
}

type manager struct {
	basepath string
	db       *gorm.DB
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
	err := ioutil.WriteFile(m.GetFullPath(path), data, 0644)
	if err != nil {
		return err
	}
	// TODO save to git and save commit id to db
	err = m.saveTime(path)
	if err != nil {
		return err
	}
	return nil
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
