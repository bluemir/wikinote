package auth

import (
	"github.com/sirupsen/logrus"
)

func (manager *Manager) Can(user *User, verb Verb, resource Resource) error {
	return manager.IsAllow(resource, verb, user)
}
func (manager *Manager) IsAllow(resource Resource, verb Verb, user *User) error {
	logrus.Tracef("user: %#v", user)
	logrus.Tracef("resource: %#v", resource)
	logrus.Tracef("verb: %s", verb)

	roles, err := manager.getBindingRoles(user)
	if err != nil {
		return err
	}

	logrus.Tracef("binding roles: %#v", roles)

	for name, role := range roles {
		if role.IsAllow(resource, verb) {
			logrus.Debugf("Allow with role '%s'", name)
			return nil
		}
	}
	if user != nil {
		return ErrForbidden
	} else {
		return ErrUnauthorized
	}
}
func (manager *Manager) getBindingRoles(user *User) (map[string]Role, error) {
	roles := map[string]struct{}{}
	x := struct{}{}

	if user != nil {
		for _, role := range manager.Binding["user/"+user.Name] {
			roles[role] = x
		}
		for group := range user.Groups {
			roles[group] = x
			for _, role := range manager.Binding["group/"+group] {
				roles[role] = x
			}
		}
	} else {
		roles[manager.Group.Unauthorized] = x
		for _, role := range manager.Binding["group/"+manager.Group.Unauthorized] {
			roles[role] = x
		}
	}

	logrus.Tracef("roles: %#v", roles)

	result := map[string]Role{}
	for name := range roles {
		if role, ok := manager.Roles[name]; ok {
			result[name] = role
		}
	}

	return result, nil
}
