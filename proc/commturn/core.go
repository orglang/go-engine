package commturn

import (
	"fmt"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/semcomm"
	"orglang/go-engine/adt/semterm"
	"orglang/go-engine/proc/termexp"
)

// aka Sem
type TurnRec interface {
	turn() identity.ADT
}

func ChnlID(rec TurnRec) identity.ADT { return rec.turn() }

type PubRec struct {
	CommRef semcomm.CommRef
	ImplRef semterm.TermRef // TODO выпилить
	ChnlID  identity.ADT
	ValExp  termexp.ExpRec
}

func (r PubRec) turn() identity.ADT { return r.ChnlID }

type SubRec struct {
	CommRef semcomm.CommRef
	ImplRef semterm.TermRef // TODO выпилить
	ChnlID  identity.ADT
	ContExp termexp.ExpRec
}

func (r SubRec) turn() identity.ADT { return r.ChnlID }

func ErrRecTypeUnexpected(got TurnRec) error {
	return fmt.Errorf("step rec unexpected: %T", got)
}

func ErrRecTypeMismatch(got, want TurnRec) error {
	return fmt.Errorf("step rec mismatch: want %T, got %T", want, got)
}
