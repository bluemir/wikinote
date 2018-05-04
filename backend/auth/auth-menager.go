package auth

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Manager interface {
	GetRules() []Rule
	PutRule(role string, action ...string)

	IsAllow(role string, actions ...string) bool
}

type Rule struct {
	Role   string
	Action string
}

const (
	defaultRule = `
admin: view, edit, user
editor: view, edit, attach
viewer: view
guest: view
	`
)

func NewManager(db *gorm.DB) (Manager, error) {
	if !db.HasTable(&Rule{}) {
		// only first time
		db.CreateTable(&Rule{})
		lines := strings.Split(defaultRule, "\n")
		for _, line := range lines {
			p := strings.SplitN(line, ":", 2)
			if len(p) < 2 {
				// skip error line
				continue
			}
			role := p[0]
			actions := strings.Split(p[1], ",")
			for _, action := range actions {
				rule := &Rule{Role: role, Action: strings.Trim(action, " ")}
				db.Where(rule).FirstOrCreate(rule)
			}
		}
	}

	return &manager{db}, nil
}

type manager struct {
	db *gorm.DB
}

func (m *manager) GetRules() []Rule {
	return nil
}
func (m *manager) PutRule(role string, allow ...string) {
	// write to file?

}

// IsAllow, actions are 'and' condition
func (m *manager) IsAllow(role string, actions ...string) bool {
	if role == "root" {
		return true
	}
	for _, action := range actions {
		logrus.Debugf("check auth %s:%s", role, action)
		rule := &Rule{Role: role, Action: action}
		c := 0
		m.db.Model(&Rule{}).Where(rule).Count(&c)
		if c == 0 {
			logrus.Debugf("check auth %s:%s %d", role, action, c)
			return false
		}
	}
	return true
}
