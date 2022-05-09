package auth

import (
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

func (m *Manager) CreateUser(username string, Labels map[string]string) error {
	salt := xid.New().String()

	if err := m.db.Create(&User{
		Name:   username,
		Labels: Labels,
		Salt:   salt,
	}).Error; err != nil {
		return errors.Wrapf(err, "User already exist")
	}

	return nil
}
func (m *Manager) GetUser(username string) (*User, bool, error) {
	user := &User{}
	if err := m.db.Where(&User{Name: username}).Take(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}
func (m *Manager) UpdateUser(user *User) error {
	if user.ID == 0 {
		return errors.Errorf("user ID not found")
	}
	return m.db.Save(user).Error
}
