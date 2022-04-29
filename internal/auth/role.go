package auth

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	yaml "gopkg.in/yaml.v3"
)

type RuleExpr struct {
	*vm.Program
}

func (rule *RuleExpr) MarshalYAML() (interface{}, error) {
	return map[string]string{"source": string(rule.Source.Content())}, nil
}
func (rule *RuleExpr) UnmarshalYAML(value *yaml.Node) error {
	source := value.Value

	p, err := expr.Compile(source, expr.AsBool())
	if err != nil {
		return err
	}

	rule.Program = p
	return nil
}
func (rule *RuleExpr) isFulfill(resource Resource, verb Verb) bool {
	env := map[string]interface{}{
		//"resource": resource.Keys(),
		"verb": verb,
	}
	vm.Run(rule.Program, env)
	return false
}

func (role *Role) IsAllow(resource Resource, verb Verb) bool {
	for _, rule := range role.Rules {
		for _, v := range rule.Verbs {
			if ResourceExprs(rule.Resource).isFulfill(resource) && v == verb {
				return true
			}
		}
	}
	return false
}
