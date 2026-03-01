package procstep

import (
	"fmt"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/procexp"
)

type CommSpec struct {
	ExecRef implsem.SemRef
	ProcES  procexp.ExpSpec
}

// aka Sem
type CommRec interface {
	step() identity.ADT
}

func ChnlID(rec CommRec) identity.ADT { return rec.step() }

type PubRec struct {
	ImplRef implsem.SemRef
	ChnlID  identity.ADT
	ValExp  procexp.ExpRec
}

func (r PubRec) step() identity.ADT { return r.ChnlID }

type SubRec struct {
	ImplRef implsem.SemRef
	ChnlID  identity.ADT
	ContExp procexp.ExpRec
}

func (r SubRec) step() identity.ADT { return r.ChnlID }

func ErrRecTypeUnexpected(got CommRec) error {
	return fmt.Errorf("step rec unexpected: %T", got)
}

func ErrRecTypeMismatch(got, want CommRec) error {
	return fmt.Errorf("step rec mismatch: want %T, got %T", want, got)
}
