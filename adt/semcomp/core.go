package semcomp

import (
	"fmt"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
)

type CompRef struct {
	CompID identity.ADT
	// revision number
	CompRN seqnum.ADT
}

func NewRef() CompRef {
	return CompRef{identity.New(), seqnum.New()}
}

func ErrConcurrentModification(ref CompRef) error {
	return fmt.Errorf("concurrent modification: %T%+v", ref, ref)
}
