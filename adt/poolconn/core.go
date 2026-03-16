package poolconn

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/option"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/seqnum"
)

type ConnRec struct {
	CommRef commsem.SemRef
}

type ConnMod struct {
	CommRef commsem.SemRef
	CommON  option.ADT[seqnum.ADT]
	Steps   []poolstep.StepRec
}

type ConnQry struct {
	CommRef commsem.SemRef
	ChnlID  option.ADT[identity.ADT]
}

// aka Configuration
type ConnSnap struct {
	CommRef commsem.SemRef
	Steps   []poolstep.StepRec
}

func (s ConnSnap) NextStep() poolstep.StepRec {
	if len(s.Steps) > 0 {
		return s.Steps[0]
	}
	return nil
}
