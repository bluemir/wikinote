package auth

import (
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type Manager struct {
	Config
	db   *gorm.DB
	salt string
}
type Config struct {
	Group struct {
		Unauthorized string
		Newcomer     []string
	}
	Roles   map[string]Role
	Binding map[string][]string
}

func New(db *gorm.DB, salt string, config *Config) (*Manager, error) {
	if err := db.AutoMigrate(
		&User{},
		&Group{},
		&Role{},
		&Token{},
	); err != nil {
		return nil, err
	}

	// dump
	buf, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	logrus.Tracef("roles: \n%s", string(buf))

	return &Manager{
		Config: *config,
		db:     db,
		salt:   salt,
	}, nil
}
