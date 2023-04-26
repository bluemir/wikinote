package auth

import (
	"gorm.io/gorm"
)

func (m *Manager) Default(name, unhashedKey string) (*User, error) {
	user := User{}

	if err := m.db.Where(&User{
		Name: name,
	}).Take(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

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

	return &user, nil
}
