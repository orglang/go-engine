package core

import (
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
	ConnRef connsem.SemRef
	CommON  option.ADT[seqnum.ADT]
	Turns   []commturn.TurnRec
}

type ConnQry struct {
	ConnRef connsem.SemRef
	ChnlID  option.ADT[identity.ADT]
}

// aka Configuration
type ConnSnap struct {
	ConnRef connsem.SemRef
	Turns   []commturn.TurnRec
}

func (s ConnSnap) NextTurn() commturn.TurnRec {
	if len(s.Turns) > 0 {
		return s.Turns[0]
	}
	return nil
}
