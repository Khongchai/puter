package parser

type Precedence int

const (
	PrecLowest = iota
	PrecIn
	PrecAssignment
	PrecLogical
	PrecEquals
	PrecLessGreater
	PrecSum
	PrecProduct
	PrecExponent
	PrecPrefix
	PrecCall
)
