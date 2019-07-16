package backend

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkgs/auth"
	"github.com/bluemir/wikinote/pkgs/fileattr"
)

type FileClause interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
	List(path string) ([]os.FileInfo, error)
	Search(query string) (interface{}, error)

	Attr(path string) fileattr.PathClause
	AttrStore() fileattr.Store
	AuthzObject(path string) auth.Object

	GetFullPath(path string) string
}

type fileClause struct {
	*backend
}

func (b *fileClause) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(b.GetFullPath(path))
}

func (b *fileClause) Write(path string, data []byte) error {
	fullpath := b.GetFullPath(path)
	dirpath := filepath.Dir(fullpath)
	if err := os.MkdirAll(dirpath, 0755); err != nil {
		return err
	}

	d, err := b.Plugin().PreSave(path, data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(b.GetFullPath(path), d, 0644)
	if err != nil {
		return err
	}

	err = b.Plugin().PostSave(path, data)
	if err != nil {
		return err
	}
	return nil
}
func (b *fileClause) List(path string) ([]os.FileInfo, error) {
	p := b.GetFullPath(path)
	ext := filepath.Ext(path)
	p = p[:len(p)-len(ext)]
	logrus.Debug(p)
	return ioutil.ReadDir(p)
}
func (b *fileClause) Search(query string) (interface{}, error) {
	//TODO query to pattern(regexp)
	return search(b.basePath, query)
}
func (b *fileClause) Attr(path string) fileattr.PathClause {
	return b.fileAttrStore.Path(path)
}
func (b *fileClause) AttrStore() fileattr.Store {
	return b.fileAttrStore
}
func (b *fileClause) AuthzObject(path string) auth.Object {
	return &authzObject{b.Attr(path)}
}
func (b *fileClause) GetFullPath(path string) string {
	return filepath.Join(b.basePath, path)
}

type authzObject struct {
	fileattr.PathClause
}

func (obj *authzObject) Attr(key string) string {
	value, _ := obj.Get(key)
	return value
}
