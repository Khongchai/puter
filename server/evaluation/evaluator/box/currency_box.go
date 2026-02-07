package box

import (
	"fmt"
)

type CurrencyBox struct {
	Number *NumberBox
	Unit   Currency
}

func (cb *CurrencyBox) Inspect() string {
	text := fmt.Sprintf("%s %s", cb.Number.Inspect(), cb.Unit)
	return text
}

func (cb *CurrencyBox) Type() BoxType {
	return CURRENCY_BOX
}

type Currency = string

var _ BinaryNumberOperatable = (*CurrencyBox)(nil)

func (cb *CurrencyBox) OperateBinary(right Box, operator BinaryOperation[float64], currencyConverter ValueConverter) (Box, error) {
	switch r := right.(type) {
	case *NumberBox:
		return &CurrencyBox{Number: NewNumberbox(operator(cb.Number.Value, r.Value), r.NumberType), Unit: cb.Unit}, nil
	case *CurrencyBox:
		if r.Unit == cb.Unit {
			return &CurrencyBox{Number: NewNumberbox(operator(cb.Number.Value, r.Number.Value), r.Number.NumberType), Unit: cb.Unit}, nil
		}

		// convert left to right
		leftConverted, err := currencyConverter(cb.Number.Value, cb.Unit, r.Unit)
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

func (nb *CurrencyBox) OperateIn(left Box, keyword string, currencyConverter ValueConverter) (Box, error) {
	if nb.Unit == keyword {
		return &CurrencyBox{Number: nb.Number, Unit: nb.Unit}, nil
	}

	isNumberKeyword, numberType := IsNumberKeyword(keyword)
	if isNumberKeyword {
		return &CurrencyBox{Number: NewNumberbox(nb.Number.Value, numberType), Unit: nb.Unit}, nil
	}

	converted, err := currencyConverter(nb.Number.Value, nb.Unit, keyword)
	if err != nil {
		return nil, err
	}
	return &CurrencyBox{Number: NewNumberbox(converted, nb.Number.NumberType), Unit: keyword}, nil
}
