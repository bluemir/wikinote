package metadata

import (
	"encoding/json"

	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type FileStoreConfig struct {
	Ext    string
	Prefix string
	Path   string
}
type FileStore struct {
	*FileStoreConfig
}

func (store *FileStore) Take(path, key string) (string, error) {
	kv, err := store.read(path)
	if err != nil {
		return "", err
	}
	return kv[key], nil
}
func (store *FileStore) Save(path, key, value string) error {
	kv, err := store.read(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	kv[key] = value

	if err := store.write(path, kv); err != nil {
		return err
	}
	return nil
}
func (store *FileStore) Delete(path, key string) error {
	kv, err := store.read(path)
	if err != nil {
		return err
	}
	delete(kv, key)

	if err := store.write(path, kv); err != nil {
		return err
	}
	return ErrNotImplemented
}
func (store *FileStore) getFullPath(path string) string {
	return filepath.Join(store.Path, path+store.Ext)
}
func (store *FileStore) read(path string) (map[string]string, error) {
	kv := map[string]string{}

	buf, err := os.ReadFile(store.getFullPath(path))
	if err != nil {
		return kv, err
	}
	if err := json.Unmarshal(buf, &kv); err != nil {
		return kv, err
	}
	return kv, nil
}
func (store *FileStore) write(path string, kv map[string]string) error {
	buf, err := json.Marshal(kv)
	if err != nil {
		return err
	}

	if err := os.WriteFile(store.getFullPath(path), buf, 0644); err != nil {
		return err
	}
	return err
}
