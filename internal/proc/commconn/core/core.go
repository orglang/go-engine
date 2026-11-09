package core

import (
	"fmt"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"

	commact "orglang/go-engine/internal/proc/commact/core"
)

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

func ErrSnapshotConflict(got seqnum.ADT) error {
	return fmt.Errorf("connection snapshot changed concurrently: offset %v", got)
}

func (s ConnSnap) NextAct() commact.ActRec {
	if len(s.Acts) > 0 {
		return s.Acts[0]
	}
	return nil
}
