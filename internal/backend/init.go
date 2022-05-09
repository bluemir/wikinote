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
	"github.com/bluemir/wikinote/internal/backend/attr"
	"github.com/bluemir/wikinote/internal/backend/files"
	"github.com/bluemir/wikinote/internal/plugins"
)

func initDB(wikipath string) (*gorm.DB, error) {
	dbPath := filepath.Join(wikipath, ".app/wikinote.db")

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

func initFileAttr(db *gorm.DB) (*attr.Store, error) {
	return attr.New(db)
}
func initAuth(db *gorm.DB, salt string, roles []auth.Role) (*auth.Manager, error) {
	return auth.New(db, salt, roles)
}

func initAdminUser(auth *auth.Manager, users map[string]string) error {
	for name, key := range users {
		if key == "" {
			key = xid.New().String()
			logrus.Warnf("generate key: '%s' '%s'", name, key)
		}
		// ensure user
		user, ok, err := auth.GetUser(name)
		if err != nil {
			return err
		}
		if !ok {
			err := auth.CreateUser(name, map[string]string{
				"role/admin": "true",
			})
			if err != nil {
				return err
			}
		} else {
			user.Labels["role/admin"] = "true"
			if err := auth.UpdateUser(user); err != nil {
				return err
			}
		}

		if err := auth.RevokeTokenAll(name); err != nil {
			return err
		}
		if _, err := auth.IssueToken(name, key, nil); err != nil {
			return err
		}
	}
	return nil
}
func initPlugins(configs []plugins.PluginConfig, attr *attr.Store) (*plugins.Manager, error) {
	return plugins.New(configs, attr)
}
func initFileStore(wikipath string) (*files.FileStore, error) {
	return files.New(wikipath)
}
