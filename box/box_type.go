package box

// Represent all objects within puter
type Box interface {
	Type() BoxType
	Inspect() string
}

type BoxType string

const (
	NUMBER_BOX       = "NUMBER"
	BOOLEAN_BOX      = "BOOLEAN"
	RETURN_VALUE_BOX = "RETURN_VALUE"
	FUNCTION_BOX     = "FUNCTION"
	BUILTIN_BOX      = "BUILTIN"
	IDENTIFIER_BOX   = "BUILTIN"
	CURRENCY_BOX     = "CURRENCY"
)
