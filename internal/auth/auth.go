package auth

import (
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type Manager struct {
	db   *gorm.DB
	salt string

	roles   map[string]Role
	binding map[string][]string
}
type Config struct {
	Roles   []Role
	Binding map[string][]string
}

func New(db *gorm.DB, salt string, config *Config) (*Manager, error) {
	if err := db.AutoMigrate(
		&User{},
		&Token{},
	); err != nil {
		return nil, err
	}

	result := map[string]Role{}
	for _, role := range config.Roles {
		result[role.Name] = role
	}

	// dump
	buf, err := yaml.Marshal(result)
	if err != nil {
		return nil, err
	}
	logrus.Tracef("roles: \n%s", string(buf))

	return &Manager{
		db:      db,
		salt:    salt,
		roles:   result,
		binding: config.Binding,
	}, nil
}
