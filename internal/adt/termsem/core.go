package termsem

import (
	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

// SemRef identifies a Term within its owning Description revision.
type SemRef struct {
	TermID identity.ADT
	DescID identity.ADT
	DescRN seqnum.ADT
}

func New(root descroot.SemRef) SemRef {
	return SemRef{TermID: identity.New(), DescID: root.DescID, DescRN: root.DescRN}
}
