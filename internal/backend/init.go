package backend

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/bluemir/wikinote/internal/backend/files"
)

func initDB(dbPath string) (*gorm.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}
	var log logger.Interface
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		log = logger.Default
	} else {
		log = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: log, // TODO make logrus adapter
	})
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
