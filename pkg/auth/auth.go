package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type Manager struct {
	db    *gorm.DB
	roles map[string]Role
}

func New(db *gorm.DB, roles []Role) (*Manager, error) {
	if err := db.AutoMigrate(
		&User{},
		&Token{},
	).Error; err != nil {
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

func (m *Manager) Subject(token *Token) (Subject, error) {
	subject := Subject{Token: token}
	if token == nil {
		return subject, nil // nil token mean guest
	}

	// find user
	user, ok, err := m.GetUser(token.UserName)
	if !ok {
		return subject, errors.Errorf("user not found")
	}
	if err != nil {
		return subject, err
	}

	subject.User = user

	return subject, nil
}
