package box

import (
	"fmt"
	"puter/evaluation/ast"
	"puter/unit"
)

type PercentBox struct {
	Value float64
}

func (nb *PercentBox) Inspect() string {
	formatted := fmt.Sprintf("%g%%", nb.Value)
	return formatted
}

func (nb *PercentBox) Type() BoxType {
	return PERCENT_BOX
}

var _ BinaryNumberOperatable = (*PercentBox)(nil)

func (pb *PercentBox) OperateBinaryNumber(right Box, operator func(a, b float64) float64, _ *unit.Converters) (Box, error) {
	switch r := right.(type) {
	case *NumberBox:
		return &NumberBox{Value: operator(r.Value, (pb.Value/100)*r.Value)}, nil
	case *CurrencyBox:
		return &CurrencyBox{Number: NewNumberbox(operator(r.Number.Value, (pb.Value/100)*r.Number.Value), r.Number.NumberType), Unit: r.Unit}, nil
	case *PercentBox:
		return &PercentBox{Value: operator(pb.Value, r.Value)}, nil
	default:
		return nil, fmt.Errorf("Cannot perform this operation on %s and %s", pb.Type(), right.Type())
	}
}

var _ BinaryBooleanOperatable = (*PercentBox)(nil)

func (left *PercentBox) OperateBinaryBoolean(right Box, operator *ast.Token, converters *unit.Converters) (Box, error) {
	r, is := right.(*PercentBox)
	if !is {
		return &BooleanBox{Value: false}, nil
	}

	result := func() bool {
		switch operator.Type {
		case ast.EQ:
			return left.Value == r.Value
		case ast.NOT_EQ:
			return left.Value != r.Value
		case ast.LT:
			return left.Value < r.Value
		case ast.GT:
			return left.Value > r.Value
		case ast.LTE:
			return left.Value <= r.Value
		case ast.GTE:
			return left.Value >= r.Value
		default:
			return false
		}
	}()

	return &BooleanBox{Value: result}, nil
}
