package metadata

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GormStoreConfig struct {
	DB *gorm.DB
}
type GormStore struct {
	*gorm.DB
}

type GormEntry struct {
	Path  string `gorm:"primary_key"`
	Key   string `gorm:"primary_key"`
	Value string
}

func (store *GormStore) Take(path, key string) (string, error) {
	entry := &GormEntry{
		Path: path,
		Key:  key,
	}
	if err := store.DB.Take(entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}
	return entry.Value, nil
}
func (store *GormStore) Save(path, key, value string) error {
	return store.DB.Save(&GormEntry{
		Path:  path,
		Key:   key,
		Value: value,
	}).Error
}
func (store *GormStore) Delete(path, key string) error {
	return store.DB.Delete(&GormEntry{
		Path: path,
		Key:  key,
	}).Error
}
