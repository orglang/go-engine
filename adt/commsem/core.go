package commsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
)

type SemRef struct {
	ChnlID identity.ADT
	ChnlON revnum.ADT
}

func (ref SemRef) Negate() SemRef {
	ref.ChnlON = -ref.ChnlON
	return ref
}
