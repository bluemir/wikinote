package files

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type FileStore struct {
	wikipath string
}

func New(ctx context.Context, wikipath string) (*FileStore, error) {
	return &FileStore{wikipath}, nil
}

func (fs *FileStore) Read(path string) ([]byte, error) {
	return os.ReadFile(fs.getFullPath(path))
}
func (fs *FileStore) ReadStream(path string) (io.ReadSeekCloser, os.FileInfo, error) {
	fullpath := fs.getFullPath(path)

	fi, err := os.Stat(fullpath)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Open(fullpath)
	if err != nil {
		return nil, nil, err
	}
	return file, fi, nil
}
func (fs *FileStore) Write(path string, data []byte) error {
	fullpath := fs.getFullPath(path)
	if err := os.MkdirAll(filepath.Dir(fullpath), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(fullpath, data, 0644); err != nil {
		return err
	}

	return nil
}
func (fs *FileStore) WriteStream(path string, r io.Reader) error {
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
func (fs *FileStore) List(path string) ([]FileInfo, error) {
	files, err := os.ReadDir(fs.getFullPath(path))
	if err != nil {
		return nil, err
	}
	ret := []FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name()[0] == '.' {
			continue
		}
		ret = append(ret, FileInfo{
			Name: file.Name(),
			Path: filepath.Join(path, file.Name()),
		})
	}
	return ret, nil
}
func (fs *FileStore) Move(oldPath, newPath string) error {
	if filepath.Ext(oldPath) != filepath.Ext(newPath) {
		return errors.New("ext not matched")
	}

	_, err := os.Stat(fs.getFullPath(newPath))
	if err != nil && err == os.ErrNotExist {
		return errors.WithStack(err)
	}
	if err == nil {
		return errors.New("target path already used.")
	}

	if err := os.MkdirAll(filepath.Dir(fs.getFullPath(newPath)), 0755); err != nil {
		return errors.WithStack(err)
	}

	if err := os.Rename(fs.getFullPath(oldPath), fs.getFullPath(newPath)); err != nil {
		return errors.WithStack(err)
	}

	msg := fmt.Sprintf("This page is moved to [%s](%s).", newPath, newPath)
	if err := os.WriteFile(fs.getFullPath(oldPath), []byte(msg), 0664); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type FileInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (fs *FileStore) getFullPath(path string) string {
	return filepath.Join(fs.wikipath, path)
}
