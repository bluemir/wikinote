package auth

import (
	"context"

	"github.com/sirupsen/logrus"
)

func (manager *Manager) IsAllow(resource Resource, verb Verb, user *User) error {
	return manager.Can(user, verb, resource)
}

func (manager *Manager) Can(user *User, verb Verb, resource Resource) error {
	logrus.Tracef("user: %#v", user)
	logrus.Tracef("resource: %#v", resource)
	logrus.Tracef("verb: %s", verb)

	roles, err := manager.getAssignedRoles(context.Background(), user)
	if err != nil {
		return err
	}

	logrus.Tracef("binding roles: %+v", roles)

	for _, role := range roles {
		if role.IsAllow(resource, verb) {
			logrus.Debugf("Allow with role '%s'", role.Name)
			return nil
		}
	}
	if user != nil {
		return ErrForbidden
	} else {
		return ErrUnauthorized
	}
}
