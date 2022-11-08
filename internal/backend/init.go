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
	"github.com/bluemir/wikinote/internal/backend/files"
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

func initAuth(db *gorm.DB, salt string, conf *auth.Config) (*auth.Manager, error) {
	// wikinote internal role: admin, it could not change
	conf.Roles["admin"] = auth.Role{
		Rules: []auth.Rule{
			{},
		},
	}

	return auth.New(db, salt, conf)
}

func initAdminUser(authm *auth.Manager, users map[string]string) error {
	for name, key := range users {
		if key == "" {
			key = xid.New().String()
			logrus.Warnf("generate key: '%s' '%s'", name, key)
		}
		// ensure user
		user, ok, err := authm.GetUser(name)
		if err != nil {
			return err
		}
		if !ok {
			err := authm.CreateUser(&auth.User{
				Name:   name,
				Groups: map[string]struct{}{"admin": {}},
			})
			if err != nil {
				return err
			}
		} else {
			user.AddGroup("admin")
			if err := authm.UpdateUser(user); err != nil {
				return err
			}
		}

		if err := authm.RevokeTokenAll(name); err != nil {
			return err
		}
		if _, err := authm.IssueToken(name, key, nil); err != nil {
			return err
		}
	}
	return nil
}

func initFileStore(wikipath string) (*files.FileStore, error) {
	return files.New(wikipath)
}
