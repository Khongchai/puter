package box

import (
	"fmt"
	"puter/lib"
)

type BooleanBox struct {
	Value lib.Promise[bool]
}

func (bb *BooleanBox) Inspect() *lib.Promise[string] {
	value := bb.Value.Await()
	return lib.NewResolvedPromise(fmt.Sprintf("%t", value))
}

func (bb *BooleanBox) Type() BoxType {
	return BOOLEAN_BOX
}
