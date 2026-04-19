package typesem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
)

type SemRef struct {
	TypeID identity.ADT
	TypeRN seqnum.ADT
}

func New() SemRef {
	return SemRef{identity.New(), seqnum.New()}
}
