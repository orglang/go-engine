package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/connsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"
	commturn "orglang/go-engine/internal/proc/commturn/core"
)

type ConnRec struct {
	ConnRef  connsem.SemRef
	OffsetNr seqnum.ADT
}

type ConnMod struct {
	ConnRef  connsem.SemRef
	OffsetNr     option.ADT[seqnum.ADT]
	PrevOffsetNr seqnum.ADT
	Turns        []commturn.TurnRec
}

type ConnQry struct {
	ConnRef connsem.SemRef
	ChnlID  option.ADT[identity.ADT]
}

// aka Configuration
type ConnSnap struct {
	ConnRef  connsem.SemRef
	OffsetNr seqnum.ADT
	Turns    []commturn.TurnRec
}

func ErrOffsetConflict(got seqnum.ADT) error {
	return fmt.Errorf("Connection offset changed concurrently: got %v", got)
}

func (s ConnSnap) NextTurn() commturn.TurnRec {
	if len(s.Turns) > 0 {
		return s.Turns[0]
	}
	return nil
}
