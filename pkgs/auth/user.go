package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func (m *Manager) CreateUser(u *User) error {
	if err := m.store.Create(u).Error; err != nil {
		return errors.Wrapf(err, "User already exist")
	}
	return nil
}
func (m *Manager) GetUser(username string) (*User, bool, error) {
	user := &User{}
	if err := m.store.Where(&User{Name: username}).First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}
func (m *Manager) ListUser(filter ...string) ([]User, error) {
	users := []User{}
	if err := m.store.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (m *Manager) DeleteUser(username string) error {
	if err := m.store.Where(&User{Name: username}).Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}
