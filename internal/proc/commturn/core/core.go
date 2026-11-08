package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/commsem"
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

// aka Sem
type TurnRec interface {
	turn() identity.ADT
}

func ChnlID(rec TurnRec) identity.ADT { return rec.turn() }

type PubRec struct {
	CommRef commsem.SemRef
	CompRef compsem.SemRef
	ChnlID  identity.ADT
	ValExp  termexp.ExpRec
}

func (r PubRec) turn() identity.ADT { return r.ChnlID }

type SubRec struct {
	CommRef commsem.SemRef
	CompRef compsem.SemRef
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
