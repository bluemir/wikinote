package auth

import (
	"strings"

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
	Op    ResourceExprOp
}

func parseResourceExpr(src string) ResourceExpr {
	if l := len(src); l < 1 || (l < 2 && src[0] == '~') {
		return ResourceExpr{}
	}
	switch src[0] {
	case '~':
		switch src[1] {
		case '+':
			return ResourceExpr{
				Value: src[2:],
				Op:    OpNotContain,
			}
		case '%':
			return ResourceExpr{
				Value: src[2:],
				Op:    OpNotIn,
			}
		default:
			return ResourceExpr{
				Value: src[1:],
				Op:    OpNotEqual,
			}
		}
	case '+':
		return ResourceExpr{
			Value: src[1:],
			Op:    OpContain,
		}
	case '^':
		return ResourceExpr{
			Value: src[1:],
			Op:    OpHasPrefix,
		}
	case '$':
		return ResourceExpr{
			Value: src[1:],
			Op:    OpHasSuffix,
		}
	case '%':
		return ResourceExpr{
			Value: src[1:],
			Op:    OpIn,
		}
	default:
		return ResourceExpr{
			Value: src,
		}
	}
}

type ResourceExprOp int

const (
	OpEqual ResourceExprOp = iota
	OpNotEqual
	OpContain
	OpNotContain
	OpHasPrefix
	OpHasSuffix
	OpIn
	OpNotIn
	//OpRegexp
)

type ResourceExprDecorator func(b bool) bool

func NotDecorator(b bool) bool {
	return !b
}

func (expr ResourceExpr) isFulfill(v string) bool {
	switch expr.Op {
	case OpNotEqual:
		return v != expr.Value
	case OpContain:
		return strings.Contains(v, expr.Value)
	case OpNotContain:
		return !strings.Contains(v, expr.Value)
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
