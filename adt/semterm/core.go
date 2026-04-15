package semterm

import (
	"fmt"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
)

type TermRef struct {
	TermID identity.ADT
	// revision number
	TermRN seqnum.ADT
}

func NewRef() TermRef {
	return TermRef{identity.New(), seqnum.New()}
}

func ErrConcurrentModification(ref TermRef) error {
	return fmt.Errorf("concurrent modification: %T%+v", ref, ref)
}
