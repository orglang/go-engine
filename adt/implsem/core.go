package implsem

import (
	"fmt"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
	"orglang/go-engine/adt/uniqsym"
)

type SemRef struct {
	ImplID identity.ADT
	// revision number
	ImplRN seqnum.ADT
}

func NewRef() SemRef {
	return SemRef{identity.New(), seqnum.New()}
}

type SemRec struct {
	ImplRef SemRef
	ImplQN  uniqsym.ADT
	Kind    semKind
}

type semKind int16

const (
	unkKind semKind = iota
	PoolKind
	ProcKind
)

func ErrConcurrentModification(ref SemRef) error {
	return fmt.Errorf("concurrent modification: %T%+v", ref, ref)
}
