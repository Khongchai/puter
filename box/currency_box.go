package box

import (
	"fmt"
)

type CurrencyBox struct {
	Number *NumberBox
	Unit   Currency
}

func (bb *CurrencyBox) Inspect() string {
	text := fmt.Sprintf("%g %s", bb.Number.Value, bb.Unit)
	return text
}

func (bb *CurrencyBox) Type() BoxType {
	return CURRENCY_BOX
}

type Currency = string
