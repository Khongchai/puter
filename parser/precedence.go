package parser

type Precedence int

const (
	PrecAssignment = iota + 1
	PrecEquals
	PrecLessGreater
	PrecSum
	PrecProduct
	PrecExponent
	PrecPrefix
	PrecCall
)
