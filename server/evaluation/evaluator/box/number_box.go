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
