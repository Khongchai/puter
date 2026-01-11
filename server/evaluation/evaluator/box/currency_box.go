package box

import (
	"fmt"
)

type CurrencyBox struct {
	Value float64
	Unit  Currency
}

func (bb *CurrencyBox) Inspect() string {
	text := fmt.Sprintf("%g %s", bb.Value, bb.Unit)
	return text
}

func (bb *CurrencyBox) Type() BoxType {
	return CURRENCY_BOX
}

type Currency = string
