package metadata

import (
	"context"

	"gorm.io/gorm"
)

type Store interface {
	//FindByPath(path string) ([]Item, error)
	//FindByKey(key string) ([]Item, error)
	Take(path, key string) (string, error)
	Save(path, key, value string) error
	Delete(path, key string) error
}

var _ Store = &GormStore{}

type Item struct {
	Path  string
	Key   string
	Value string
}

type Config struct {
}

func New(ctx context.Context, db *gorm.DB) (Store, error) {
	if err := db.AutoMigrate(&GormEntry{}); err != nil {
		return nil, err
	}
	return &GormStore{db}, nil
}
