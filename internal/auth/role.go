package auth

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	yaml "gopkg.in/yaml.v3"
)

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

func (role *Role) IsAllow(resource Resource, verb Verb) bool {
	for _, rule := range role.Rules {
		for _, v := range rule.Verbs {
			if rule.Resource.isSubsetOf(resource) && v == verb {
				return true
			}
		}
	}
	return false
}
