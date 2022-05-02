package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseResourceExpr(t *testing.T) {
	tcs := []struct {
		Source string
		Result ResourceExpr
	}{
		{
			Source: "value",
			Result: ResourceExpr{
				Value: "value",
				Op:    OpEqual,
			},
		},
		{
			Source: "~value",
			Result: ResourceExpr{
				Value: "value",
				Op:    OpNotEqual,
			},
		},
		{
			Source: "+value",
			Result: ResourceExpr{
				Value: "value",
				Op:    OpContain,
			},
		},
		{
			Source: "~+value",
			Result: ResourceExpr{
				Value: "value",
				Op:    OpNotContain,
			},
		},
		{
			Source: "^value",
			Result: ResourceExpr{
				Value: "value",
				Op:    OpHasPrefix,
			},
		},
		{
			Source: "$value",
			Result: ResourceExpr{
				Value: "value",
				Op:    OpHasSuffix,
			},
		},
		{
			Source: "%value1,value2",
			Result: ResourceExpr{
				Value: "value1,value2",
				Op:    OpIn,
			},
		},
		{
			Source: "~%value1,value2",
			Result: ResourceExpr{
				Value: "value1,value2",
				Op:    OpNotIn,
			},
		},
		{
			Source: "",
			Result: ResourceExpr{
				Value: "",
			},
		},
	}
	for _, tc := range tcs {
		expr := parseResourceExpr(tc.Source)

		assert.Equal(t, expr, tc.Result)
	}
}
func TestIsFulfilled(t *testing.T) {

	tcs := []struct {
		Expr   ResourceExpr
		Value  string
		Result bool
	}{
		{
			Expr: ResourceExpr{
				Value: "value",
				Op:    OpEqual,
			},
			Value:  "value",
			Result: true,
		},
		{
			Expr: ResourceExpr{
				Value: "value",
				Op:    OpNotEqual,
			},
			Value:  "value",
			Result: false,
		},
		{
			Expr: ResourceExpr{
				Value: "value",
				Op:    OpNotEqual,
			},
			Value:  "value1",
			Result: true,
		},
		{
			Expr: ResourceExpr{
				Value: "value",
				Op:    OpContain,
			},
			Value:  "value1",
			Result: true,
		},
	}
	for _, tc := range tcs {
		res := tc.Expr.isFulfill(tc.Value)

		assert.Equal(t, res, tc.Result)
	}
}
