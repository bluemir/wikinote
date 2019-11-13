package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func (m *Manager) CreateUser(username string, Labels map[string]string) error {
	if err := m.db.Create(&User{
		Name:   username,
		Labels: Labels,
	}).Error; err != nil {
		return errors.Wrapf(err, "User already exist")
	}
	return nil
}
func (m *Manager) GetUser(username string) (*User, bool, error) {
	user := &User{}
	if err := m.db.Where(&User{Name: username}).Take(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
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
