package parser

type Precedence int

const (
	PrecLowest = iota
	PrecAssignment
	PrecLogical
	PrecEquals
	PrecLessGreater
	PrecBitwiseOperators
	PrecSum
	PrecProduct
	PrecExponent
	PrecPrefix
	PrecIn
	PrecCall
	PrecPercent
)
