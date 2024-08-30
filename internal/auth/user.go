package auth

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type User struct {
	Name     string    `gorm:"primary_key" json:"name"`
	CreateAt time.Time `json:"createAt"`
	Groups   []Group   `gorm:"many2many:members;" json:"groups"`
	Labels   Labels    `gorm:"type:bytes;serializer:gob" json:"labels" expr:"labels"`
	Salt     string    `json:"-"`
}

func (u User) Subject() Subject {
	return Subject{
		Kind: KindUser,
		Name: u.Name,
	}
}

func (m *Manager) CreateUser(ctx context.Context, user *User) (*User, error) {
	// overwrite salt
	user.Salt = xid.New().String()
	user.CreateAt = time.Now()

	if len(user.Groups) == 0 {
		for _, group := range m.conf.Group.NewUserGroups {
			user.Groups = append(user.Groups, Group{
				Name: group,
			})
		}
	}

	if err := m.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, errors.Wrapf(err, "User already exist")
	}

	return user, nil
}
func (m *Manager) GetUser(ctx context.Context, username string) (*User, bool, error) {
	user := &User{}
	if err := m.db.WithContext(ctx).Preload("Groups").Where(&User{Name: username}).Take(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, errors.WithStack(err)
	}
	return user, true, nil
}
func (m *Manager) ListUsers(ctx context.Context) ([]User, error) {
	users := []User{}
	if err := m.db.WithContext(ctx).Preload("Groups").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (m *Manager) UpdateUser(ctx context.Context, user *User) error {
	logrus.Tracef("update user: %+v", user)

	return errors.WithStack(m.db.WithContext(ctx).Save(user).Error)
}

func (m *Manager) DeleteUser(ctx context.Context, name string) error {
	return m.db.WithContext(ctx).Delete(User{
		Name: name,
	}).Error
}

func (user *User) AddGroup(group string) {
	user.Groups = append(user.Groups, Group{Name: group})
}
func (user *User) RemoveGroup(group string) {
	for i, g := range user.Groups {
		if g.Name == group {
			user.Groups = append(user.Groups[:i], user.Groups[i+1:]...)
			return
		}
	}
}
