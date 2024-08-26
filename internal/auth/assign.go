package auth

import (
	"context"
	"errors"

	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

type Assign struct {
	Subject Subject `gorm:"embedded"`
	Roles   Set     `gorm:"type:bytes;serializer:gob"`
}
type Subject struct {
	Kind Kind   `gorm:"primaryKey;size:256" expr:"kind"`
	Name string `gorm:"primaryKey;size:256" expr:"name"`
}

type Kind string

const (
	KindUser  Kind = "user"
	KindGroup Kind = "group"
	KindGuest Kind = "guest"
)

func (m *Manager) AssignRole(ctx context.Context, subject Subject, roles ...string) error {
	assign := &Assign{
		Subject: subject,
		Roles:   Set{},
	}
	if err := m.db.WithContext(ctx).FirstOrCreate(&assign).Error; err != nil {
		return err
	}
	assign.Roles.Add(roles...)

	if err := m.db.WithContext(ctx).Save(assign).Error; err != nil {
		return err
	}

	return nil
}

func (m *Manager) DiscardRole(ctx context.Context, subject Subject, roles ...string) error {
	return ErrNotImplements
}

func (m *Manager) getAssignedRoles(ctx context.Context, user *User) ([]Role, error) {
	if user == nil { // it is guest
		assign, err := m.getAssign(ctx, Subject{
			Kind: KindGuest,
		})
		if err != nil {
			return nil, err
		}

		return m.findRoles(ctx, assign.Roles)
	}

	roleNames := Set{}

	assign, err := m.getAssign(ctx, Subject{
		Kind: KindUser,
		Name: user.Name,
	})
	if err != nil {
		return nil, err
	}
	maps.Copy(roleNames, assign.Roles)

	for group := range user.Groups {
		assign, err := m.getAssign(ctx, Subject{
			Kind: KindGroup,
			Name: group,
		})
		if err != nil {
			return nil, err
		}
		maps.Copy(roleNames, assign.Roles)

		roleNames.Add(group)
	}

	return m.findRoles(ctx, roleNames)
}
func (m *Manager) getAssign(ctx context.Context, subject Subject) (*Assign, error) {
	assign := Assign{}
	if err := m.db.WithContext(ctx).Where(&Assign{
		Subject: subject,
	}).Take(&assign).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	return &assign, nil
}

func (m *Manager) findRoles(ctx context.Context, roleNames Set) ([]Role, error) {
	roles := []Role{}
	for roleName := range roleNames {
		role := Role{
			Name: roleName,
		}
		if err := m.db.WithContext(ctx).Find(&role).Error; err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}
	return roles, nil
}
