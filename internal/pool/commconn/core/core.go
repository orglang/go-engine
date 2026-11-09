package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"

	commact "orglang/go-engine/internal/pool/commact/core"
)

type ConnRec struct {
	ConnID          identity.ADT
	ProviderCompRef compsem.SemRef
	ClientCompRef   compsem.SemRef
	OffsetNr        seqnum.ADT
	ActNr           seqnum.ADT
}

type ConnMod struct {
	ConnID         identity.ADT
	ProviderCompID identity.ADT
	ProviderCompRN seqnum.ADT
	OffsetNr       option.ADT[seqnum.ADT]
	Acts           []commact.ActRec
}

type ConnQry struct {
	ConnID identity.ADT
	ChnlID option.ADT[identity.ADT]
}

// aka Configuration
type ConnSnap struct {
	ConnID         identity.ADT
	ProviderCompID identity.ADT
	ProviderCompRN seqnum.ADT
	OffsetNr       seqnum.ADT
	Acts           []commact.ActRec
}

func ErrOffsetConflict(got seqnum.ADT) error {
	return fmt.Errorf("connection offset changed concurrently: got %v", got)
}

func ErrMissingConn(id identity.ADT) error {
	return fmt.Errorf("connection missing: %v", id)
}

func (s ConnSnap) NextAct() commact.ActRec {
	if len(s.Acts) > 0 {
		return s.Acts[0]
	}
	return nil
}
