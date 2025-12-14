package evaluator

import "fmt"

type ObjectType string

const (
	NUMBER_OBJ       = "NUMBER"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	IDENTIFIER_OBJ   = "BUILTIN"
)

// Represent all objects within puter
type Box interface {
	Type() ObjectType
	Inspect() string
}

type BooleanBox struct {
	value bool
}

func (bb *BooleanBox) Inspect() string {
	return fmt.Sprintf("%t", bb.value)
}

func (bb *BooleanBox) Type() ObjectType {
	return BOOLEAN_OBJ
}

type IdentBox struct {
	value string
}

func (ib *IdentBox) Inspect() string {
	return ib.value
}

func (ib *IdentBox) Type() ObjectType {
	return IDENTIFIER_OBJ
}

type NumberBox struct {
	value float64
}

func (nb *NumberBox) Inspect() string {
	return fmt.Sprintf("%g", nb.value)
}

func (nb *NumberBox) Type() ObjectType {
	return NUMBER_OBJ
}
