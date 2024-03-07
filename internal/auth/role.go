package auth

import (
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
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
