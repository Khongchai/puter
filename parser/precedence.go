package parser

type Precedence int

const (
	PrecLowest = iota
	PrecAssignment
	PrecLogical
	PrecEquals
	PrecLessGreater
	PrecSum
	PrecProduct
	PrecExponent
	PrecPrefix
	PrecIn
	PrecCall
)
