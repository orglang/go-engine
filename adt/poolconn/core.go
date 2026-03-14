package poolconn

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/option"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/seqnum"
	"orglang/go-engine/adt/uniqsym"
)

type ConnRec struct {
	CommRef commsem.SemRef
}

type ConnMod struct {
	CommRef commsem.SemRef
	CommON  option.ADT[seqnum.ADT]
	Steps   []poolstep.StepRec
}

type ConnQuery struct {
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

type Repo interface {
	InsertRec(db.Source, ConnRec) error
	UpdateRec(db.Source, ConnMod) error
	SelectRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error)
	SelectSnapByQry(db.Source, ConnQuery) (ConnSnap, error)
}
