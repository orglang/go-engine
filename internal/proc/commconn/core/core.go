package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"

	commturn "orglang/go-engine/internal/proc/commturn/core"
)

type ConnRec struct {
	ConnID          identity.ADT
	ProviderCompRef compsem.SemRef
	ClientCompRef   compsem.SemRef
	OffsetNr        seqnum.ADT
	ActNr           seqnum.ADT
}

type ConnMod struct {
	ConnID          identity.ADT
	ProviderCompRef compsem.SemRef
	ClientCompRef   compsem.SemRef
	OffsetNr        option.ADT[seqnum.ADT]
	PrevOffsetNr seqnum.ADT
	Turns        []commturn.TurnRec
}

type ConnQry struct {
	ConnID identity.ADT
	ChnlID option.ADT[identity.ADT]
}

// aka Configuration
type ConnSnap struct {
	ConnID   identity.ADT
	OffsetNr seqnum.ADT
	Turns    []commturn.TurnRec
}

func ErrOffsetConflict(got seqnum.ADT) error {
	return fmt.Errorf("connection offset changed concurrently: got %v", got)
}

func ErrSnapshotConflict(got seqnum.ADT) error {
	return fmt.Errorf("connection snapshot changed concurrently: offset %v", got)
}

func (s ConnSnap) NextTurn() commturn.TurnRec {
	if len(s.Turns) > 0 {
		return s.Turns[0]
	}
	return nil
}
