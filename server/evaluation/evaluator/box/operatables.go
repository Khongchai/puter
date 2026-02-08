package box

import "puter/unit"

type BinaryOperation[T any] = func(a, b T) T

type BinaryNumberOperatable interface {
	OperateBinary(right Box, operation BinaryOperation[float64], converters *unit.Converters) (Box, error)
}

// Can apply "in" keyword
// 20 in usd << number "in" usd
type InPrefixOperatable interface {
	OperateIn(keyword string, converters *unit.Converters) (Box, error)
}
