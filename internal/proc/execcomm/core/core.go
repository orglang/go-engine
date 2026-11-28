package core

import (
	"fmt"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"

	commact "orglang/go-engine/internal/proc/commact/core"
)

type CommMod struct {
	CommID         identity.ADT
	ProviderExecID identity.ADT
	ProviderExecRN seqnum.ADT
	OffsetNr       option.ADT[seqnum.ADT]
	Acts           []commact.ActRec
}

type CommQry struct {
	CommID identity.ADT
	ChnlID option.ADT[identity.ADT]
}

// aka Configuration
type CommSnap struct {
	CommID         identity.ADT
	ProviderExecID identity.ADT
	ProviderExecRN seqnum.ADT
	OffsetNr       seqnum.ADT
	Acts           []commact.ActRec
}

func ErrOffsetConflict(got seqnum.ADT) error {
	return fmt.Errorf("communication offset changed concurrently: got %v", got)
}

func ErrSnapshotConflict(got seqnum.ADT) error {
	return fmt.Errorf("communication snapshot changed concurrently: offset %v", got)
}

func (s CommSnap) NextAct() commact.ActRec {
	if len(s.Acts) > 0 {
		return s.Acts[0]
	}
	return nil
}
