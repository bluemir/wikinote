package auth

import (
	"context"
	"encoding/gob"
	"regexp"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&Rule{})
}

type Role struct {
	Name  string `yaml:"name"`
	Rules []Rule `yaml:"rules" gorm:"type:bytes;serializer:gob"`
}
type Rule struct {
	Verbs      []Verb          `yaml:"verbs"`
	Resources  []ResourceMatch `yaml:"resources"`
	Conditions []Condition     `yaml:"conditions"` // TODO make serializer
}
type ResourceMatch map[string]string

func (m *Manager) CreateRole(ctx context.Context, name string, rules []Rule) (*Role, error) {
	role := &Role{
		Name:  name,
		Rules: rules,
	}
	if err := m.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}
func (m *Manager) GetRole(ctx context.Context, name string) (*Role, error) {
	role := &Role{}
	if err := m.db.WithContext(ctx).Where(Role{Name: name}).Take(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}
func (m *Manager) ListRoles(ctx context.Context) ([]Role, error) {
	roles := []Role{}
	if err := m.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
func (m *Manager) UpdateRole(ctx context.Context, role *Role) error {
	if err := m.db.Where(Role{Name: role.Name}).Save(role).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
func (m *Manager) DeleteRole(ctx context.Context, name string) error {
	if err := m.db.WithContext(ctx).Where(&Role{
		Name: name,
	}).Delete(Role{
		Name: name,
	}).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (role *Role) IsAllow(resource Resource, verb Verb) bool {
	for _, rule := range role.Rules {
		if rule.isMatch(verb, resource) {
			return true
		}
	}
	return false
}
func (rule *Rule) isMatch(verb Verb, resource Resource) bool {
	return rule.isMatchVerb(verb) && rule.isMatchResource(resource) && rule.isMatchCondition(verb, resource)
}
func (rule *Rule) isMatchVerb(verb Verb) bool {
	if len(rule.Verbs) == 0 {
		return true
	}
	for _, v := range rule.Verbs {
		if verb == v {
			return true
		}
	}
	return false
}
func (rule *Rule) isMatchResource(resource Resource) bool {
	if len(rule.Resources) == 0 {
		return true
	}

	for _, match := range rule.Resources {
		if match.isMatch(resource) {
			return true
		}
	}
	return false
}
func (match ResourceMatch) isMatch(resource Resource) bool {
	for k, v := range match {
		re, err := regexp.Compile(v)
		if err != nil {
			logrus.Warn(err)
			return false
		}
		if !re.MatchString(resource.Get(k)) {
			return false
		}
	}
	return true
}

func (rule *Rule) isMatchCondition(verb Verb, resource Resource) bool {
	if len(rule.Conditions) == 0 {
		return true
	}

	for _, cond := range rule.Conditions {
		if r, err := cond.IsMatched(Context{
			Verb:     verb,
			Resource: resource,
		}); err != nil {
			logrus.Error(err)
			return false
		} else if !r {
			return false
		}
	}
	return true
}
