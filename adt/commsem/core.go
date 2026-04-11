package commsem

import (
	"fmt"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
)

type SemRef struct {
	CommID identity.ADT
	// revision number
	CommRN seqnum.ADT
}

func NewRef() SemRef {
	return SemRef{identity.New(), seqnum.New()}
}

type SemRec struct {
	CommRef SemRef
	Kind    semKind
}

type semKind int16

const (
	unkSem semKind = iota
	Pool
	Proc
)

func ErrConcurrentModification(ref SemRef) error {
	return fmt.Errorf("concurrent modification: %T%+v", ref, ref)
}
