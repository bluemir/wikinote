package metadata

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Store struct {
	*gorm.DB
}

func (store *Store) Take(ctx context.Context, path, key string) (string, error) {
	entry := &StoreItem{
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
func (store *Store) Save(ctx context.Context, path, key, value string) error {
	return store.DB.Save(&StoreItem{
		Path:  path,
		Key:   key,
		Value: value,
	}).Error
}
func (store *Store) Delete(ctx context.Context, path, key string) error {
	return store.DB.Delete(&StoreItem{
		Path: path,
		Key:  key,
	}).Error
}
func (store *Store) FindByLabels(ctx context.Context, labels map[string]string) ([]StoreItem, error) {

	cond := [][]any{}
	for k, v := range labels {
		cond = append(cond, []any{k, v})
	}

	entries := []StoreItem{}

	if err := store.DB.WithContext(ctx).
		Model(&StoreItem{}).
		Where("(key, value) in ?", cond).
		Group("path").
		Having("count(*) = ?", len(labels)).
		Find(&entries).Error; err != nil {

		return nil, err
	}

	return entries, nil
}

func (store *Store) List(ctx context.Context) ([]StoreItem, error) {
	items := []StoreItem{}
	if err := store.DB.WithContext(ctx).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
