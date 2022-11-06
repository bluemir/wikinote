package auth

import (
	"github.com/sirupsen/logrus"
)

func (manager *Manager) IsAllow(resource Resource, verb Verb, user *User) error {
	logrus.Tracef("user: %#v", user)
	logrus.Tracef("resource: %#v", resource)
	logrus.Tracef("verb: %s", verb)

	roles, err := manager.getBindingRoles(user)
	if err != nil {
		return err
	}
	logrus.Tracef("%+v", roles)
	for _, role := range roles {
		if role.IsAllow(resource, verb) {
			return nil
		}
	}
	if user != nil {
		return ErrForbidden
	} else {
		return ErrUnauthorized
	}
}
func (manager *Manager) getBindingRoles(user *User) ([]Role, error) {
	// if user nil, it's guest
	if user == nil {
		return []Role{manager.roles["guest"]}, nil
	}

	roles := user.Roles

	result := []Role{}
	for _, name := range roles {
		result = append(result, manager.roles[name])
	}

	return result, nil
}
