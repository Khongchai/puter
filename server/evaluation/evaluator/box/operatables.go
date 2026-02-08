package box

import (
	"puter/evaluation/ast"
	"puter/unit"
)

type BinaryNumberOperatable interface {
	OperateBinaryNumber(right Box, operation func(a, b float64) float64, converters *unit.Converters) (Box, error)
}

// Can apply "in" keyword
// 20 in usd << number "in" usd
type InPrefixOperatable interface {
	OperateIn(keyword string, converters *unit.Converters) (Box, error)
}

type BinaryBooleanOperatable interface {
	OperateBinaryBoolean(right Box, operator *ast.Token, converters *unit.Converters) (Box, error)
}

type HoldsNumber interface {
	GetNumber() float64
}
