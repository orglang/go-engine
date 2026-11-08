package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/connsem"
	"orglang/go-engine/internal/adt/identity"
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

// aka Sem
type TurnRec interface {
	turn() identity.ADT
}

func ChnlID(rec TurnRec) identity.ADT { return rec.turn() }

type PubRec struct {
	ConnRef connsem.SemRef
	CompRef compsem.SemRef
	ChnlID  identity.ADT
	ValExp  termexp.ExpRec
}

func (r PubRec) turn() identity.ADT { return r.ChnlID }

func (r PubRec) GetConnRef() connsem.SemRef { return r.ConnRef }

type SubRec struct {
	ConnRef connsem.SemRef
	CompRef compsem.SemRef
	ChnlID  identity.ADT
	ContExp termexp.ExpRec
}

func (r SubRec) turn() identity.ADT { return r.ChnlID }

func (r SubRec) GetConnRef() connsem.SemRef { return r.ConnRef }

func ErrRecTypeUnexpected(got TurnRec) error {
	return fmt.Errorf("step rec unexpected: %T", got)
}

func ErrRecTypeMismatch(got, want TurnRec) error {
	return fmt.Errorf("step rec mismatch: want %T, got %T", want, got)
}
