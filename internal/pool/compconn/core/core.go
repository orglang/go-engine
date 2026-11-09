package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"

	connact "orglang/go-engine/internal/pool/connact/core"
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
	Acts           []connact.ActRec
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
	Acts           []connact.ActRec
}

func ErrOffsetConflict(got seqnum.ADT) error {
	return fmt.Errorf("connection offset changed concurrently: got %v", got)
}

func ErrMissingConn(id identity.ADT) error {
	return fmt.Errorf("connection missing: %v", id)
}

func (s ConnSnap) NextAct() connact.ActRec {
	if len(s.Acts) > 0 {
		return s.Acts[0]
	}
	return nil
}
