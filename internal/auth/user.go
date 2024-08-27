package auth

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type User struct {
	Name   string `gorm:"primary_key" json:"name"`
	Groups Set    `gorm:"type:bytes;serializer:gob" json:"groups"`
	Labels Labels `gorm:"type:bytes;serializer:gob" json:"labels" expr:"labels"`
	Salt   string `json:"-"`
}

func (m *Manager) CreateUser(ctx context.Context, user *User) (*User, error) {

	// overwrite salt
	user.Salt = xid.New().String()

	if len(user.Groups) == 0 {
		user.Groups = map[string]struct{}{}
		for _, group := range m.conf.Group.Newcomer {
			user.Groups.Add(group)
		}
	}

	if err := m.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, errors.Wrapf(err, "User already exist")
	}

	return user, nil
}
func (m *Manager) GetUser(ctx context.Context, username string) (*User, bool, error) {
	user := &User{}
	if err := m.db.WithContext(ctx).Where(&User{Name: username}).Take(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, errors.WithStack(err)
	}
	return user, true, nil
}
func (m *Manager) ListUsers(ctx context.Context) ([]User, error) {
	users := []User{}
	if err := m.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (m *Manager) UpdateUser(ctx context.Context, user *User) error {

	return m.db.WithContext(ctx).Save(user).Error
}

func (m *Manager) DeleteUser(ctx context.Context, name string) error {
	return m.db.WithContext(ctx).Delete(User{
		Name: name,
	}).Error
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
