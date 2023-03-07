package auth

import (
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

func (m *Manager) CreateUser(user *User) error {
	if user.ID != 0 {
		return errors.Errorf("user ID already exist")
	}

	// overwrite salt
	user.Salt = xid.New().String()

	if len(user.Groups) == 0 {
		user.Groups = map[string]struct{}{}
		for _, group := range m.Group.Newcomer {
			user.Groups[group] = struct{}{}
		}
	}

	if err := m.db.Create(user).Error; err != nil {
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

type User struct {
	ID     uint   `gorm:"primary_key" json:"-"`
	Name   string `gorm:"unique" json:"name"`
	Groups Set    `sql:"type:json" json:"groups"`
	Labels Labels `sql:"type:json" json:"labels"`
	Salt   string `json:"-"`
}

func (user *User) AddGroup(group string) {
	user.Groups[group] = struct{}{}
}
func (user *User) RemoveGroup(group string) {
	delete(user.Groups, group)
}
func (m *Manager) GetMember(group string) ([]User, error) {
	users := []User{}
	if err := m.db.Find(&users).Error; err != nil {
		return users, nil
	}
	ret := []User{}
	for _, u := range users {
		for g := range u.Groups {
			if g == group {
				ret = append(ret, u)
				break
			}
		}
	}
	return ret, nil
}
