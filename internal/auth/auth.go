package auth

import (
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type Manager struct {
	db    *gorm.DB
	roles map[string]Role
}

func New(db *gorm.DB, roles []Role) (*Manager, error) {
	if err := db.AutoMigrate(
		&User{},
		&Token{},
		&RoleBinding{},
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

	return &Manager{db, result}, nil
}
