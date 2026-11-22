package typesem

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

type SemRef struct {
	TypeID identity.ADT
	DescID identity.ADT
	DescRN seqnum.ADT
}

func New() SemRef {
	return SemRef{TypeID: identity.New(), DescID: identity.New(), DescRN: seqnum.New()}
}

type typeKind int16

const (
	unkKind typeKind = iota
	PoolKind
	ProcKind
)
