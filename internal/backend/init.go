package backend

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/fileattr"
	"github.com/bluemir/wikinote/internal/plugins"
)

func (backend *Backend) initDB() error {
	dbPath := filepath.Join(backend.Config.Wikipath, ".app/wikinote.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return err
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to connect database")
	}
	if rawDB, err := db.DB(); err != nil {
		return err
	} else {
		rawDB.SetMaxOpenConns(1)
	}
	return nil
}

func (backend *Backend) initFileAttr() error {
	if backend.db == nil {
		return errors.Errorf("require db")
	}
	fa, err := fileattr.New(backend.db)
	if err != nil {
		return err
	}
	backend.FileAttr = fa
	return nil
}
func (backend *Backend) initAuth() error {
	if backend.db == nil {
		return errors.Errorf("require db")
	}

	m, err := auth.New(backend.db, backend.Config.File.Roles)
	if err != nil {
		return err
	}
	backend.Auth = m

	return nil
}
func (backend *Backend) initAdminUser() error {
	if backend.Auth == nil {
		return errors.Errorf("require auth")
	}

	for name, key := range backend.Config.AdminUsers {
		if key == "" {
			key = xid.New().String()
			logrus.Warnf("generate key: '%s' '%s'", name, key)
		}
		if err := backend.Auth.EnsureUser(name, map[string]string{
			"role/root": "true",
		}); err != nil {
			return err
		}
		if err := backend.Auth.RevokeTokenAll(name); err != nil {
			return err
		}
		if _, err := backend.Auth.IssueToken(name, key); err != nil {
			return err
		}
	}
	return nil
}
func (backend *Backend) initPlugins() error {
	if backend.FileAttr == nil {
		return errors.Errorf("require auth")
	}
	m, err := plugins.New(backend.Config.File.Plugins, backend.FileAttr)
	if err != nil {
		return err
	}
	backend.Plugin = m

	return nil
}
