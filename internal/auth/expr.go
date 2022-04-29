package auth

import (
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type ResourceExprs map[string]ResourceExpr

func (exprs ResourceExprs) isFulfill(resource Resource) bool {
	for k, expr := range exprs {
		if v := resource.Get(k); !expr.isFulfill(v) {
			return false
		}
	}
	return true
}

type ResourceExpr struct {
	Value string
	Op    ResourceOp
}
type ResourceOp int

const (
	Equal ResourceOp = iota
	NotEqual
	Contain
	NotContain
	HasPrefix
	HasSuffix
	Regexp
	In
)

func (expr ResourceExpr) isFulfill(v string) bool {
	switch expr.Op {
	default:
		return v == expr.Value
	}
}
func (expr ResourceExpr) MarshalYAML() (interface{}, error) {
	return expr.Value, nil
}
func (expr *ResourceExpr) UnmarshalYAML(value *yaml.Node) error {
	logrus.Tracef("%+v", value)
	str := value.Value

	expr.Value = str

	return nil
}
