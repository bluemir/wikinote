package files

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileStore struct {
	wikipath string
}

func New(wikipath string) (*FileStore, error) {
	return &FileStore{wikipath}, nil
}

func (fs *FileStore) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(fs.getFullPath(path))
}
func (fs *FileStore) ReadStream(path string) (io.ReadSeekCloser, error) {
	fullpath := fs.getFullPath(path)

	return os.Open(fullpath)
}
func (fs *FileStore) Write(path string, data []byte) error {
	fullpath := fs.getFullPath(path)
	if err := os.MkdirAll(filepath.Dir(fullpath), 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(fullpath, data, 0644); err != nil {
		return err
	}

	return nil
}
func (fs *FileStore) WriteStream(path string, r io.ReadCloser) error {
	fullpath := fs.getFullPath(path)

	if err := os.MkdirAll(filepath.Dir(fullpath), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return err
	}

	return nil
}

func (fs *FileStore) Delete(path string) error {
	return os.Remove(fs.getFullPath(path))
}

func (fs *FileStore) getFullPath(path string) string {
	return filepath.Join(fs.wikipath, path)
}
