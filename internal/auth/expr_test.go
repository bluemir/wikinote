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
				Op:    OpEqual,
				Flag:  FlagNot,
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
				Op:    OpContain,
				Flag:  FlagNot,
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
				Op:    OpIn,
				Flag:  FlagNot,
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
				Op:    OpEqual,
				Flag:  FlagNot,
			},
			Value:  "value",
			Result: false,
		},
		{
			Expr: ResourceExpr{
				Value: "value",
				Op:    OpEqual,
				Flag:  FlagNot,
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
		{
			Expr: ResourceExpr{
				Value: "value",
				Op:    OpHasSuffix,
			},
			Value:  "test-value",
			Result: true,
		},
		{
			Expr: ResourceExpr{
				Value: "value1,value2",
				Op:    OpIn,
			},
			Value:  "value2",
			Result: true,
		},
	}
	for _, tc := range tcs {
		res := tc.Expr.isFulfill(tc.Value)

		assert.Equal(t, res, tc.Result)
	}
}
func TestResourceExprToString(t *testing.T) {
	tcs := []struct {
		Result string
		Source ResourceExpr
	}{
		{
			Source: ResourceExpr{
				Value: "value",
				Op:    OpEqual,
			},
			Result: "value",
		},
		{
			Result: "~value",
			Source: ResourceExpr{
				Value: "value",
				Op:    OpEqual,
				Flag:  FlagNot,
			},
		},
		{
			Result: "+value",
			Source: ResourceExpr{
				Value: "value",
				Op:    OpContain,
			},
		},
		{
			Result: "~+value",
			Source: ResourceExpr{
				Value: "value",
				Op:    OpContain,
				Flag:  FlagNot,
			},
		},
		{
			Result: "^value",
			Source: ResourceExpr{
				Value: "value",
				Op:    OpHasPrefix,
			},
		},
		{
			Result: "$value",
			Source: ResourceExpr{
				Value: "value",
				Op:    OpHasSuffix,
			},
		},
		{
			Result: "%value1,value2",
			Source: ResourceExpr{
				Value: "value1,value2",
				Op:    OpIn,
			},
		},
		{
			Result: "~%value1,value2",
			Source: ResourceExpr{
				Value: "value1,value2",
				Op:    OpIn,
				Flag:  FlagNot,
			},
		},
		{
			Result: "",
			Source: ResourceExpr{
				Value: "",
			},
		},
	}
	for _, tc := range tcs {
		assert.Equal(t, tc.Result, tc.Source.String())
	}
}
