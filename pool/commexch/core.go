package commexch

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/option"
	"orglang/go-engine/adt/seqnum"

	"orglang/go-engine/pool/commturn"
)

type ExchRec struct {
	CommRef  commsem.SemRef
	OffsetNr seqnum.ADT
}

type ExchMod struct {
	CommRef  commsem.SemRef
	OffsetNr option.ADT[seqnum.ADT]
	Turns    []commturn.TurnRec
}

type ExchQry struct {
	CommRef commsem.SemRef
	ChnlID  option.ADT[identity.ADT]
}

// aka Configuration
type ExchSnap struct {
	CommRef commsem.SemRef
	Turns   []commturn.TurnRec
}

func (s ExchSnap) NextTurn() commturn.TurnRec {
	if len(s.Turns) > 0 {
		return s.Turns[0]
	}
	return nil
}
