package box

import (
	"fmt"
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

func IsNumberKeyword(keyword string) (bool, NumberType) {
	lowercased := strings.ToLower(keyword)
	is := lowercased == string(Decimal) || keyword == string(Binary) || keyword == string(Hex)
	if is {
		return true, NumberType(lowercased)
	}
	return false, NaN
}

var _ BinaryNumberOperatables = (*NumberBox)(nil)

func (nb *NumberBox) OperateBinary(right Box, operator BinaryOperation[float64], valueConverter ValueConverter) (Box, error) {
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
