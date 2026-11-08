package connsem

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

type SemRef struct {
	ConnID identity.ADT
	// revision number
	ConnRN seqnum.ADT
}

func New() SemRef {
	return SemRef{identity.New(), seqnum.New()}
}
