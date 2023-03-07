package events

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewTestHub() (IHub[struct{}], error) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		return nil, err
	}
	return NewHub[struct{}](db)
}
