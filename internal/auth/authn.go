package auth

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (m *Manager) Default(name, unhashedKey string) (*User, error) {
	user := User{}

	if err := m.db.Preload("Groups").Where(&User{
		Name: name,
	}).Take(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	// If the user has a bcrypt hashed password, use it
	if user.HashedPassword != "" {
		if !bcryptCompare(user.HashedPassword, unhashedKey) {
			return nil, ErrUnauthorized
		}
		return &user, nil
	}

	// Legacy fallback: check Token table with SHA-512
	token := &Token{}
	if err := m.db.Where(&Token{
		Username:  name,
		HashedKey: hash(unhashedKey, salt(name)),
	}).Take(token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	// Auto-migrate: store bcrypt hash for future logins
	hashed, err := bcryptHash(unhashedKey)
	if err != nil {
		logrus.WithError(err).Warn("failed to migrate password to bcrypt")
	} else {
		user.HashedPassword = hashed
		if err := m.db.Save(&user).Error; err != nil {
			logrus.WithError(err).Warn("failed to save migrated password")
		}
	}

	return &user, nil
}
