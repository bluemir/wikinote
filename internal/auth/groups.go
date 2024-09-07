package auth

import (
	"context"

	"github.com/pkg/errors"
)

type Group struct {
	Name    string `gorm:"primary_key" json:"name"`
	Members []User `gorm:"many2many:members;"`
}

func (g Group) Subject() Subject {
	return Subject{
		Kind: KindGroup,
		Name: g.Name,
	}
}
func (m *Manager) CreateGroup(ctx context.Context, name string) error {
	if err := m.db.WithContext(ctx).Create(&Group{
		Name: name,
	}).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *Manager) ListGroups(ctx context.Context) ([]Group, error) {
	groups := []Group{}
	if err := m.db.WithContext(ctx).Find(&groups).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return groups, nil
}
func (m *Manager) GetGroup(ctx context.Context, name string) (*Group, error) {
	group := &Group{}
	if err := m.db.WithContext(ctx).Preload("Members").Where(Group{Name: name}).Take(group).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return group, nil
}
