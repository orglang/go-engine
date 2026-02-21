package procstep

import (
	"fmt"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/procexp"
	"orglang/go-engine/adt/revnum"
)

type StepSpec struct {
	// ExecRef uniqref.ADT
	ExecRef implsem.SemRef
	ProcES  procexp.ExpSpec
}

// aka Sem
type StepRec interface {
	step() identity.ADT
}

func ChnlID(rec StepRec) identity.ADT { return rec.step() }

type MsgRec struct {
	PoolRN revnum.ADT
	// ExecRef uniqref.ADT
	ExecRef implsem.SemRef
	ChnlID  identity.ADT
	ValER   procexp.ExpRec
}

func (r MsgRec) step() identity.ADT { return r.ChnlID }

type SvcRec struct {
	PoolRN revnum.ADT
	// ExecRef uniqref.ADT
	ExecRef implsem.SemRef
	ChnlID  identity.ADT
	ContER  procexp.ExpRec
}

func (r SvcRec) step() identity.ADT { return r.ChnlID }

func ErrRecTypeUnexpected(got StepRec) error {
	return fmt.Errorf("step rec unexpected: %T", got)
}

func ErrRecTypeMismatch(got, want StepRec) error {
	return fmt.Errorf("step rec mismatch: want %T, got %T", want, got)
}
