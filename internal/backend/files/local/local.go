package local

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type LocalFileSystem struct {
	wikipath string
}

func New() *LocalFileSystem {
	return &LocalFileSystem{}
}

func (fs *LocalFileSystem) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(fs.getFullPath(path))
}

func (fs *LocalFileSystem) Write(path string, data []byte) error {
	fullpath := fs.getFullPath(path)
	dirpath := filepath.Dir(fullpath)
	if err := os.MkdirAll(dirpath, 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(fs.getFullPath(path), data, 0644); err != nil {
		return err
	}

	return nil
}

func (fs *LocalFileSystem) Delete(path string) error {
	return os.Remove(fs.getFullPath(path))
}

func (fs *LocalFileSystem) getFullPath(path string) string {
	return filepath.Join(fs.wikipath, path)
}
