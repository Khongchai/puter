package parser

type Precedence int

const (
	PrecLowest = iota
	PrecIn
	PrecAssignment
	PrecEquals
	PrecLessGreater
	PrecSum
	PrecProduct
	PrecExponent
	PrecPrefix
	PrecCall
)
