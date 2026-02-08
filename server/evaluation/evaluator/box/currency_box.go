package box

import (
	"fmt"
	"puter/evaluation/ast"
	"puter/unit"
)

type CurrencyBox struct {
	Number *NumberBox
	Unit   unit.Currency
}

func (cb *CurrencyBox) Inspect() string {
	text := fmt.Sprintf("%s %s", cb.Number.Inspect(), cb.Unit)
	return text
}

func (cb *CurrencyBox) Type() BoxType {
	return CURRENCY_BOX
}

var _ BinaryNumberOperatable = (*CurrencyBox)(nil)

func (cb *CurrencyBox) OperateBinaryNumber(right Box, operator func(a, b float64) float64, converters *unit.Converters) (Box, error) {
	switch r := right.(type) {
	case *NumberBox:
		return &CurrencyBox{Number: NewNumberbox(operator(cb.Number.Value, r.Value), r.NumberType), Unit: cb.Unit}, nil
	case *CurrencyBox:
		if r.Unit == cb.Unit {
			return &CurrencyBox{Number: NewNumberbox(operator(cb.Number.Value, r.Number.Value), r.Number.NumberType), Unit: cb.Unit}, nil
		}

		// convert left to right
		leftConverted, err := converters.ConvertCurrency(cb.Number.Value, cb.Unit, r.Unit)
		if err != nil {
			return nil, err
		}

		return &CurrencyBox{Number: NewNumberbox(operator(leftConverted, r.Number.Value), r.Number.NumberType), Unit: r.Unit}, nil
	case *PercentBox:
		// 2 + 2% = 2 + (2/200 * 2)
		return &CurrencyBox{Number: NewNumberbox(operator(cb.Number.Value, (r.Value/100)*cb.Number.Value), cb.Number.NumberType), Unit: cb.Unit}, nil
	default:
		return nil, fmt.Errorf("Cannot perform this operation on %s and %s", cb.Type(), right.Type())
	}
}

var _ InPrefixOperatable = (*CurrencyBox)(nil)

func (nb *CurrencyBox) OperateIn(keyword string, converters *unit.Converters) (Box, error) {
	if nb.Unit == keyword {
		return &CurrencyBox{Number: nb.Number, Unit: nb.Unit}, nil
	}

	isNumberKeyword, numberType := IsNumberKeyword(keyword)
	if isNumberKeyword {
		return &CurrencyBox{Number: NewNumberbox(nb.Number.Value, numberType), Unit: nb.Unit}, nil
	}

	converted, err := converters.ConvertCurrency(nb.Number.Value, nb.Unit, keyword)
	if err != nil {
		return nil, err
	}
	return &CurrencyBox{Number: NewNumberbox(converted, nb.Number.NumberType), Unit: keyword}, nil
}

var _ BinaryBooleanOperatable = (*CurrencyBox)(nil)

func (left *CurrencyBox) OperateBinaryBoolean(right Box, operator *ast.Token, converters *unit.Converters) (Box, error) {
	r, is := right.(*CurrencyBox)
	if !is {
		return &BooleanBox{Value: false}, nil
	}

	var leftAsRight, err = converters.ConvertCurrency(left.Number.Value, left.Unit, r.Unit)
	if err != nil {
		return nil, err
	}

	result := func() bool {
		switch operator.Type {
		case ast.EQ:
			return leftAsRight == r.Number.Value
		case ast.NOT_EQ:
			return leftAsRight != r.Number.Value
		case ast.LT:
			return leftAsRight < r.Number.Value
		case ast.GT:
			return leftAsRight > r.Number.Value
		case ast.LTE:
			return leftAsRight <= r.Number.Value
		case ast.GTE:
			return leftAsRight >= r.Number.Value
		default:
			return false
		}
	}()

	return &BooleanBox{Value: result}, nil
}

var _ HoldsNumber = (*CurrencyBox)(nil)

func (c *CurrencyBox) GetNumber() float64 {
	return c.Number.Value
}
