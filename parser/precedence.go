package parser

type Precedence int

const (
	PrecLowest = iota
	PrecAssignment
	PrecIn
	PrecLogical
	PrecEquals
	PrecLessGreater
	PrecSum
	PrecProduct
	PrecExponent
	PrecPrefix
	PrecCall
)
