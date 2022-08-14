package auth

import (
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type Manager struct {
	db          *gorm.DB
	salt        string
	roles       map[string]Role
	defaultRole string
}

func New(db *gorm.DB, salt string, roles []Role, defaultRole string) (*Manager, error) {
	if err := db.AutoMigrate(
		&User{},
		&Token{},
	); err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		logrus.Warn("config dosen't have role. using defualt role.")
	}

	result := map[string]Role{}
	for _, role := range roles {
		result[role.Name] = role
	}

	// dump
	buf, err := yaml.Marshal(result)
	if err != nil {
		return nil, err
	}
	logrus.Tracef("roles: \n%s", string(buf))

	return &Manager{
		db:    db,
		salt:  salt,
		roles: result,
	}, nil
}
