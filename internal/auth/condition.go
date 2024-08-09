package auth

import (
	"github.com/expr-lang/expr"
)

type Condition string

func (cond Condition) IsMatched(ctx Context) (bool, error) {
	p, err := expr.Compile(string(cond), expr.AsBool())
	if err != nil {
		return false, err
	}
	result, err := expr.Run(p, ctx)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

type Context struct {
	//User     *User    `expr:"user"`
	//Subject  Subject  `expr:"subject"`
	Verb     Verb     `expr:"verb"`
	Resource Resource `expr:"resource"`
}
