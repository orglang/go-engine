package commsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
)

type SemRef struct {
	CommID identity.ADT
	// revision number
	CommRN seqnum.ADT
}

func New() SemRef {
	return SemRef{identity.New(), seqnum.New()}
}
