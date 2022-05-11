package attr

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Attribute struct {
	Path  string `gorm:"primary_key"`
	Key   string `gorm:"primary_key"`
	Value string
}
type ListOption struct {
	Order string
	Limit int
}

type Store struct {
	db *gorm.DB
}

type FindOption func()
type AttrStore interface {
	Find(q *Attribute, opts ...FindOption) ([]Attribute, error)
	Take(path, key string) (*Attribute, error)
	Save(attr *Attribute) error
	Delete(attr *Attribute) error
}

func New(db *gorm.DB) (*Store, error) {
	if err := db.AutoMigrate(
		&Attribute{},
	); err != nil {
		return nil, errors.Wrap(err, "auto migrate is failed")
	}

	return &Store{db}, nil
}

func (store *Store) Find(attr *Attribute) ([]Attribute, error) {
	attrs := []Attribute{}
	if err := store.db.Where(attr).Find(&attrs).Error; err != nil {
		return nil, err
	}
	return attrs, nil
}
func (store *Store) Search(attr *Attribute, opt *ListOption) ([]Attribute, error) {
	attrs := []Attribute{}
	err := store.db.Where(attr).Order(opt.Order).Limit(opt.Limit).Find(&attrs).Error
	return attrs, err
}
func (store *Store) Take(attr *Attribute) (*Attribute, error) {
	result := &Attribute{}
	if err := store.db.Where(attr).Take(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (store *Store) Save(attr *Attribute) error {
	//TODO check empty value
	return store.db.Save(attr).Error
}
func (store *Store) Delete(attr *Attribute) error {
	return store.db.Where(attr).Delete(attr).Error
}
func IsNotFound(err error) bool {
	return gorm.ErrRecordNotFound == err
}

// need Raw method ?
