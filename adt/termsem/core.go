package termsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
)

type SemRef struct {
	TermID identity.ADT
	// revision number
	TermRN seqnum.ADT
}

func New() SemRef {
	return SemRef{identity.New(), seqnum.New()}
}
