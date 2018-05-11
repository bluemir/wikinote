package backend

import "github.com/sirupsen/logrus"

type AuthClause interface {
	GetRules() []Rule
	PutRule(role string, action ...string)

	IsAllow(role string, actions ...string) bool
}

type authClause backend

func (b *authClause) GetRules() []Rule {
	return nil
}
func (b *authClause) PutRule(role string, action ...string) {

}

func (b *authClause) IsAllow(role string, actions ...string) bool {
	if role == "root" {
		return true
	}
	for _, action := range actions {
		logrus.Debugf("check auth %s:%s", role, action)
		rule := &Rule{Role: role, Action: action}
		c := 0
		b.db.Model(&Rule{}).Where(rule).Count(&c)
		if c == 0 {
			logrus.Debugf("check auth %s:%s %d", role, action, c)
			return false
		}
	}
	return true
}

type Rule struct {
	Role   string
	Action string
}
