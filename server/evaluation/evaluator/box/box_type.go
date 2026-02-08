package box

// Represent all objects within puter
type Box interface {
	Type() BoxType
	Inspect() string
}

type BoxType string

const (
	PERCENT_BOX      = "PERCENT"
	NUMBER_BOX       = "NUMBER"
	BOOLEAN_BOX      = "BOOLEAN"
	RETURN_VALUE_BOX = "RETURN_VALUE"
	FUNCTION_BOX     = "FUNCTION"
	BUILTIN_BOX      = "BUILTIN"
	CURRENCY_BOX     = "CURRENCY"
	FIXED_UNIT_BOX   = "FIXED_UNIT_BOX" // cm, km, lbs, pounds, kb, gb, etc.
)
