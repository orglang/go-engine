package commexch

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/option"
	"orglang/go-engine/adt/semcomm"
	"orglang/go-engine/adt/seqnum"
	"orglang/go-engine/proc/commturn"
)

type ExchRec struct {
	CommRef  semcomm.CommRef
	OffsetNr seqnum.ADT
}

type ExchMod struct {
	CommRef semcomm.CommRef
	CommON  option.ADT[seqnum.ADT]
	Turns   []commturn.TurnRec
}

type ExchQry struct {
	CommRef semcomm.CommRef
	ChnlID  option.ADT[identity.ADT]
}

// aka Configuration
type ExchSnap struct {
	CommRef semcomm.CommRef
	Turns   []commturn.TurnRec
}

func (s ExchSnap) NextTurn() commturn.TurnRec {
	if len(s.Turns) > 0 {
		return s.Turns[0]
	}
	return nil
}
