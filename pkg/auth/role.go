package auth

import (
	"io/ioutil"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

const defaultRoles = `
- name: root
  rules:
  - actions:
- name: default
  rules:
  - actions:
    - wiki:read
`

type Role struct {
	Name  string `yaml:"name"`
	Rules []Rule `yaml:"rules"`
}
type Rule struct {
	Objects []ObjectRule `yaml:"objects"`
	Actions []string     `yaml:"actions"`
}
type ObjectRule struct {
	*vm.Program
}

func (rule *ObjectRule) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(rule.Source.Content())
}
func (rule *ObjectRule) UnmarshalYAML(value *yaml.Node) error {
	source := value.Value

	p, err := expr.Compile(source)
	if err != nil {
		return err
	}

	rule.Program = p
	return nil
}

func loadRole(roleFile string) (map[string]Role, error) {
	log := logrus.WithField("method", "auth.loadRole")

	result := map[string]Role{}

	roles := []Role{}

	if err := yaml.Unmarshal([]byte(defaultRoles), &roles); err != nil {
		log.Warnf("'%s' not parsed. %s", roleFile, err)
		return nil, err
	}

	for _, role := range roles {
		result[role.Name] = role
	}

	log.Tracef("default: \n%s\n", must(yaml.Marshal(result)))

	buf, err := ioutil.ReadFile(roleFile)
	if err != nil {
		log.Warnf("'%s' not exist. %s", roleFile, err)
		buf = []byte(``)
	}
	if err = yaml.Unmarshal(buf, &roles); err != nil {
		log.Warnf("'%s' not parsed. %s", roleFile, err)
		return nil, err
	}

	for _, role := range roles {
		result[role.Name] = role
	}

	log.Tracef("all: \n%s\n", must(yaml.Marshal(result)))

	return result, nil
}
func must(buf []byte, err error) string {
	if err != nil {
		panic(buf)
	}
	return string(buf)
}

func (rule *Rule) IsMatchedAction(action string) bool {
	if len(rule.Actions) == 0 {
		return true // match all
	}
	for _, a := range rule.Actions {
		if a == action {
			return true
		}
	}
	return false
}
func (rule *Rule) IsNotMatchedAction(action string) bool {
	return !rule.IsMatchedAction(action)
}
func (rule *Rule) IsMatchedObject(obj map[string]string) bool {
	if len(rule.Objects) == 0 {
		return true // match all
	}
	for _, o := range rule.Objects {
		output, err := expr.Run(o.Program, obj)
		if err != nil {
			continue
		}
		r, ok := output.(bool)
		if !ok {
			continue
		}
		if r {
			return true
		}
	}
	return false
}

func (rule *Rule) IsNotMatchedObject(obj map[string]string) bool {
	return !rule.IsMatchedObject(obj)
}
