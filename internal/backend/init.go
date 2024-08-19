package backend

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/backend/files"
)

func initDB(dbPath string) (*gorm.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}
	if rawDB, err := db.DB(); err != nil {
		return nil, err
	} else {
		rawDB.SetMaxOpenConns(1)
	}

	return db, nil
}

func initFileStore(wikipath string) (*files.FileStore, error) {
	return files.New(wikipath)
}
