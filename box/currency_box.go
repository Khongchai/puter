package box

type Currency = string

type CurrencyBox struct {
	number *NumberBox
	unit   Currency
}
