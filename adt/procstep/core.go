package procstep

import (
	"fmt"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/procexp"
)

type StepSpec struct {
	ImplRef implsem.SemRef
	ProcExp procexp.ExpSpec
}

// aka Sem
type StepRec interface {
	step() identity.ADT
}

func ChnlID(rec StepRec) identity.ADT { return rec.step() }

type PubRec struct {
	CommRef commsem.SemRef
	ImplRef implsem.SemRef // TODO выпилить
	ChnlID  identity.ADT
	ValExp  procexp.ExpRec
}

func (r PubRec) step() identity.ADT { return r.ChnlID }

type SubRec struct {
	CommRef commsem.SemRef
	ImplRef implsem.SemRef // TODO выпилить
	ChnlID  identity.ADT
	ContExp procexp.ExpRec
}

func (r SubRec) step() identity.ADT { return r.ChnlID }

func ErrRecTypeUnexpected(got StepRec) error {
	return fmt.Errorf("step rec unexpected: %T", got)
}

func ErrRecTypeMismatch(got, want StepRec) error {
	return fmt.Errorf("step rec mismatch: want %T, got %T", want, got)
}
