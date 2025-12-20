package box

type IdentBox struct {
	Identifier string
	Value      string
}

func (ib *IdentBox) Inspect() string {
	return ib.Value
}

func (ib *IdentBox) Type() BoxType {
	return IDENTIFIER_BOX
}
