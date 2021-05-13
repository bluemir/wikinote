package auth

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	yaml "gopkg.in/yaml.v3"
)

type Role struct {
	Name  string `yaml:"name"`
	Rules []Rule `yaml:"rules"`
}
type Rule struct {
	Objects []*ObjectRule `yaml:"objects"`
	Actions []string      `yaml:"actions"`
}
type ObjectRule struct {
	*vm.Program
}

func (rule *ObjectRule) MarshalYAML() (interface{}, error) {
	return map[string]string{"source": string(rule.Source.Content())}, nil
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
