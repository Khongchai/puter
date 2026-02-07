package box

type BinaryOperation[T any] = func(a, b T) T

type BinaryNumberOperatable interface {
	OperateBinary(right Box, operation BinaryOperation[float64], currencyConverter ValueConverter) (Box, error)
}

// Can apply "in" keyword
// 20 in usd << number "in" usd
type InPrefixOperatable interface {
	OperateIn(keyword string, currencyConverter ValueConverter) (Box, error)
}
