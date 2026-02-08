package box

import (
	"fmt"
	"puter/evaluation/ast"
	"puter/unit"
	"strings"
)

type NumberType string

const (
	NaN     NumberType = "NaN"
	Decimal NumberType = "decimal"
	Binary  NumberType = "binary"
	Hex     NumberType = "hex"
)

type NumberBox struct {
	Value      float64
	NumberType NumberType
}

func NewNumberbox(value float64, numberType NumberType) *NumberBox {
	return &NumberBox{Value: value, NumberType: numberType}
}

func (nb *NumberBox) Inspect() string {
	formatted := func() string {
		if nb.NumberType == NaN {
			return "NaN"
		}
		if nb.NumberType == Binary {
			return fmt.Sprintf("0b%b", int(nb.Value))
		}
		if nb.NumberType == Hex {
			return fmt.Sprintf("0x%x", int(nb.Value))
		}
		return fmt.Sprintf("%g", nb.Value)
	}()
	return formatted
}

func (nb *NumberBox) Type() BoxType {
	return NUMBER_BOX
}

var _ BinaryNumberOperatable = (*NumberBox)(nil)

func (nb *NumberBox) OperateBinaryNumber(right Box, operator func(a, b float64) float64, converters *unit.Converters) (Box, error) {
	switch r := right.(type) {
	case *NumberBox:
		return &NumberBox{Value: operator(nb.Value, r.Value)}, nil
	case *CurrencyBox:
		return &CurrencyBox{Number: NewNumberbox(operator(nb.Value, r.Number.Value), r.Number.NumberType), Unit: r.Unit}, nil
	case *PercentBox:
		return &NumberBox{Value: operator(nb.Value, (r.Value/100)*nb.Value)}, nil
	default:
		return nil, fmt.Errorf("Cannot perform this operation on %s and %s", nb.Type(), right.Type())
	}
}

var _ InPrefixOperatable = (*NumberBox)(nil)

func (nb *NumberBox) OperateIn(keyword string, converters *unit.Converters) (Box, error) {
	isNumberKeyword, numberType := IsNumberKeyword(keyword)
	if isNumberKeyword {
		return NewNumberbox(nb.Value, numberType), nil
	}
	isFixedUnitKeyword, fixedUnitType := unit.IsFixedUnitKeyword(keyword)
	if isFixedUnitKeyword {
		return NewFixedUnitBox(NewNumberbox(nb.Value, nb.NumberType), fixedUnitType), nil
	}
	return &CurrencyBox{
		Number: NewNumberbox(nb.Value, nb.NumberType),
		Unit:   keyword,
	}, nil
}

func IsNumberKeyword(keyword string) (bool, NumberType) {
	lowercased := strings.ToLower(keyword)
	is := lowercased == string(Decimal) || keyword == string(Binary) || keyword == string(Hex)
	if is {
		return true, NumberType(lowercased)
	}
	return false, NaN
}

var _ BinaryBooleanOperatable = (*NumberBox)(nil)

func (left *NumberBox) OperateBinaryBoolean(right Box, operator *ast.Token, converters *unit.Converters) (Box, error) {
	r, is := right.(*NumberBox)
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

var _ HoldsNumber = (*NumberBox)(nil)

func (n *NumberBox) GetNumber() float64 {
	return n.Value
}
