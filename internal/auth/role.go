package auth

func (m *Manager) ListRoles() ([]Role, error) {
	roles := []Role{}
	if err := m.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
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
	return rule.isMatchVerb(verb) && rule.isMatchResource(resource)
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
		if !v.MatchString(resource.Get(k)) {
			return false
		}
	}
	return true
}
