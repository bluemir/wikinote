package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Manager struct {
	db    *gorm.DB
	roles map[string]Role
}

func New(db *gorm.DB, roleFile string) (*Manager, error) {
	// Init db
	if err := db.AutoMigrate(
		&User{},
		&Token{},
	).Error; err != nil {
		return nil, err
	}

	// Load role
	roles, err := loadRole(roleFile)
	if err != nil {
		return nil, err
	}

	return &Manager{db, roles}, nil
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
