package box

type BinaryOperation[T any] = func(a, b T) T

type BinaryNumberOperatables interface {
	OperateBinary(right Box, operation BinaryOperation[float64], valueConverter ValueConverter) (Box, error)
}
