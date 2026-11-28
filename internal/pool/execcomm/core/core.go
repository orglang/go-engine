package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"

	commact "orglang/go-engine/internal/pool/commact/core"
)

type CommRec struct {
	CommID          identity.ADT
	ProviderRootRef execroot.RootRef
	ClientRootRef   execroot.RootRef
	OffsetNr        seqnum.ADT
	ActNr           seqnum.ADT
}

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
	return fmt.Errorf("connection offset changed concurrently: got %v", got)
}

func ErrMissingConn(id identity.ADT) error {
	return fmt.Errorf("connection missing: %v", id)
}

func (s CommSnap) NextAct() commact.ActRec {
	if len(s.Acts) > 0 {
		return s.Acts[0]
	}
	return nil
}
