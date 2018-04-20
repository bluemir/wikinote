package auth

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type Manager interface {
	GetRules() []Rule
	PutRule(role string, action ...string)

	IsAllow(role string, actions ...string) bool
}

func NewManager() (Manager, error) {

	return &manager{
		parseRule(),
	}, nil
}
func parseRule() map[string][]string {
	result := map[string][]string{}
	lines := strings.Split(defaultRule, "\n")
	for _, line := range lines {
		p := strings.SplitN(line, ":", 2)
		if len(p) < 2 {
			// skip error line
			continue
		}
		role := p[0]
		rules := strings.Split(p[1], ",")
		for _, rule := range rules {
			result[role] = append(result[role], strings.Trim(rule, " \t\n\r"))
		}
	}

	logrus.Debugf("parse default rules: %+v", result)
	return result
}

const (
	defaultRule = `
admin: view, edit, user
editor: view, edit, attach
viewer: view
guest: view
	`
)

type Rule struct {
	role    string
	allowed []string
}
type manager struct {
	rules map[string][]string
}

func (m *manager) GetRules() []Rule {
	return nil
}
func (m *manager) PutRule(role string, allow ...string) {
	// write to file?

}
func (m *manager) IsAllow(role string, actions ...string) bool {
	if role == "root" {
		return true
	}
	for _, action := range actions {
		if !m.isAllow(role, action) {
			return false
		}
	}
	return true
}
func (m *manager) isAllow(role, action string) bool {
	rule, ok := m.rules[role]
	logrus.Debugf("isAllow: %s - %s %+v", role, action, rule)
	if ok {
		for _, a := range rule {
			if action == a {
				return true
			}
		}
	}
	return false
}
