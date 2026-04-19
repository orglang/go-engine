package compsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
)

type SemRef struct {
	CompID identity.ADT
	// revision number
	CompRN seqnum.ADT
}

func New() SemRef {
	return SemRef{identity.New(), seqnum.New()}
}
