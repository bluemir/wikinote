package auth

import (
	"github.com/sirupsen/logrus"
)

func (manager *Manager) IsAllow(resource Resource, verb Verb, user *User) (bool, error) {
	logrus.Tracef("user: %#v", user)
	logrus.Tracef("resource: %#v", resource)
	logrus.Tracef("verb: %s", verb)

	roles, err := manager.getBindingRoles(user)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.IsAllow(resource, verb) {
			return true, nil
		}
	}

	return false, nil
}
func (manager *Manager) getBindingRoles(user *User) ([]Role, error) {
	// if user nil, it's guest
	if user == nil {
		return []Role{manager.roles["guest"]}, nil
	}

	bindings := []RoleBinding{}
	if err := manager.db.Where(RoleBinding{
		Username: user.Name,
	}).Find(bindings).Error; err != nil {
		return nil, err
	}

	result := []Role{}

	for _, b := range bindings {
		result = append(result, manager.roles[b.Rolename])
	}

	return result, nil
}
