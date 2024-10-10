package metadata

import (
	"context"

	"gorm.io/gorm"
)

type IStore interface {
	//FindByPath(path string) ([]Item, error)
	//FindByKey(key string) ([]Item, error)
	Take(ctx context.Context, path, key string) (string, error)
	Save(ctx context.Context, path, key, value string) error
	Delete(ctx context.Context, path, key string) error

	List(ctx context.Context) ([]StoreItem, error)

	FindByLabels(ctx context.Context, labels map[string]string) ([]StoreItem, error)
}

var _ IStore = (*Store)(nil)

type StoreItem struct {
	Path  string `gorm:"primary_key"`
	Key   string `gorm:"primary_key"`
	Value string
}

func New(ctx context.Context, db *gorm.DB) (*Store, error) {
	if err := db.AutoMigrate(&StoreItem{}); err != nil {
		return nil, err
	}
	return &Store{db}, nil
}
