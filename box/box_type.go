package box

import (
	"puter/ast"
	"puter/lib"
)

// Represent all objects within puter
type Box interface {
	Type() BoxType
	Inspect() *lib.Promise[string]
	TokenValue() *ast.Token
}

type BoxType string

const (
	NUMBER_BOX       = "NUMBER"
	BOOLEAN_BOX      = "BOOLEAN"
	RETURN_VALUE_BOX = "RETURN_VALUE"
	FUNCTION_BOX     = "FUNCTION"
	BUILTIN_BOX      = "BUILTIN"
	IDENTIFIER_BOX   = "BUILTIN"
)
