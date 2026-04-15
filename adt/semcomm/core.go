package semcomm

import (
	"fmt"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
)

type CommRef struct {
	CommID identity.ADT
	// revision number
	CommRN seqnum.ADT
}

func NewRef() CommRef {
	return CommRef{identity.New(), seqnum.New()}
}

type SemRec struct {
	CommRef CommRef
	Kind    semKind
}

type semKind int16

const (
	unkKind semKind = iota
	PoolKind
	ProcKind
)

func ErrConcurrentModification(ref CommRef) error {
	return fmt.Errorf("concurrent modification: %T%+v", ref, ref)
}
