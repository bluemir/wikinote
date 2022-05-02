package auth

import (
	"strings"

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
	Flag  ResourceExprFlag
}

func parseResourceExpr(src string) ResourceExpr {
	if len(src) == 0 {
		return ResourceExpr{}
	}

	flag := FlagNomal
	if src[0] == '~' {
		flag = FlagNot
		src = src[1:]
	}

	if len(src) == 0 {
		return ResourceExpr{
			Flag: flag,
		}
	}

	switch src[0] {
	case '+':
		return ResourceExpr{
			Value: src[1:],
			Op:    OpContain,
			Flag:  flag,
		}
	case '^':
		return ResourceExpr{
			Value: src[1:],
			Op:    OpHasPrefix,
			Flag:  flag,
		}
	case '$':
		return ResourceExpr{
			Value: src[1:],
			Op:    OpHasSuffix,
			Flag:  flag,
		}
	case '%':
		return ResourceExpr{
			Value: src[1:],
			Op:    OpIn,
			Flag:  flag,
		}
	default:
		return ResourceExpr{
			Value: src,
			Flag:  flag,
		}
	}
}

type ResourceExprOp int

const (
	OpEqual ResourceExprOp = iota
	OpContain
	OpHasPrefix
	OpHasSuffix
	OpIn
)

type ResourceExprFlag int

const (
	FlagNomal ResourceExprFlag = iota
	FlagNot
)

type ResourceExprDecorator func(b bool) bool

func NotDecorator(b bool) bool {
	return !b
}

func (expr ResourceExpr) isFulfill(v string) bool {

	if expr.Flag == FlagNot {
		return !(ResourceExpr{
			Op:    expr.Op,
			Value: expr.Value,
		}).isFulfill(v)
	}

	switch expr.Op {
	case OpContain:
		return strings.Contains(v, expr.Value)
	case OpHasPrefix:
		return strings.HasPrefix(v, expr.Value)
	case OpHasSuffix:
		return strings.HasSuffix(v, expr.Value)
	case OpIn:
		arr := strings.Split(expr.Value, ",")
		for _, str := range arr {
			if str == v {
				return true
			}
		}
		return false
	default:
		return v == expr.Value
	}
}
func (expr *ResourceExpr) String() string {
	str := ""
	switch expr.Flag {
	case FlagNot:
		str += "~"
	}
	switch expr.Op {
	case OpContain:
		str += "+"
	case OpHasPrefix:
		str += "^"
	case OpHasSuffix:
		str += "$"
	case OpIn:
		str += "%"
	}
	return str + expr.Value
}
func (expr ResourceExpr) MarshalYAML() (interface{}, error) {
	return expr.String(), nil
}
func (expr *ResourceExpr) UnmarshalYAML(value *yaml.Node) error {
	e := parseResourceExpr(value.Value)

	expr.Value = e.Value
	expr.Op = e.Op

	return nil
}
