package core

import (
	"orglang/go-engine/internal/adt/commsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"

	commturn "orglang/go-engine/internal/pool/commturn/core"
)

type ConnRec struct {
	CommRef  commsem.SemRef
	OffsetNr seqnum.ADT
}

type ConnMod struct {
	CommRef  commsem.SemRef
	OffsetNr option.ADT[seqnum.ADT]
	Turns    []commturn.TurnRec
}

type ConnQry struct {
	CommRef commsem.SemRef
	ChnlID  option.ADT[identity.ADT]
}

// aka Configuration
type ConnSnap struct {
	CommRef commsem.SemRef
	Turns   []commturn.TurnRec
}

func (s ConnSnap) NextTurn() commturn.TurnRec {
	if len(s.Turns) > 0 {
		return s.Turns[0]
	}
	return nil
}
