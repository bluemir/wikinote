package fileattr

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type FileAttr struct {
	Path  string `gorm:"primary_key"`
	Key   string `gorm:"primary_key"`
	Value string
}

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*Store, error) {
	if err := db.AutoMigrate(
		&FileAttr{},
	).Error; err != nil {
		return nil, errors.Wrap(err, "auto migrate is failed")
	}

	return &Store{db}, nil
}

func (store *Store) Find(attr *FileAttr) ([]FileAttr, error) {
	attrs := []FileAttr{}
	if err := store.db.Where(attr).Find(&attrs).Error; err != nil {
		return nil, err
	}
	return attrs, nil
}
func (store *Store) Take(attr *FileAttr) (*FileAttr, error) {
	result := &FileAttr{}
	if err := store.db.Where(attr).Take(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (store *Store) Save(attr *FileAttr) error {
	//TODO check empty value
	return store.db.Save(attr).Error
}

// need Raw method ?
