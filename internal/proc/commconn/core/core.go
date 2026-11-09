package core

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"

	commturn "orglang/go-engine/internal/proc/commturn/core"
)

type ConnMod struct {
	ConnID         identity.ADT
	ProviderCompID identity.ADT
	ProviderCompRN seqnum.ADT
	OffsetNr       option.ADT[seqnum.ADT]
	Turns          []commturn.TurnRec
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
	Turns          []commturn.TurnRec
}

func (s ConnSnap) NextTurn() commturn.TurnRec {
	if len(s.Turns) > 0 {
		return s.Turns[0]
	}
	return nil
}
