package poolcomm

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/option"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/seqnum"
)

type ConnRec struct {
	CommRef commsem.SemRef
	// offset number
	CommON seqnum.ADT
}

type CommMod struct {
	CommRef commsem.SemRef
	CommON  option.ADT[seqnum.ADT]
	Steps   []poolstep.StepRec
}

type CommQry struct {
	CommRef commsem.SemRef
	ChnlID  option.ADT[identity.ADT]
}

// aka Configuration
type CommSnap struct {
	CommRef commsem.SemRef
	Steps   []poolstep.StepRec
}

func (s CommSnap) NextStep() poolstep.StepRec {
	if len(s.Steps) > 0 {
		return s.Steps[0]
	}
	return nil
}
