package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

// step (aka Act)
type TurnRec interface {
	turn()
	GetConnID() identity.ADT
	GetActNr() seqnum.ADT
}

// publication (aka msg)
type PubRec struct {
	ConnID  identity.ADT
	ActNr   seqnum.ADT
	CompRef compsem.SemRef
	ChnlID  identity.ADT
	ValExp  termexp.ExpRec
}

func (r PubRec) turn() {}
func (r PubRec) GetConnID() identity.ADT { return r.ConnID }
func (r PubRec) GetActNr() seqnum.ADT { return r.ActNr }

// subscription (aka srv)
type SubRec struct {
	ConnID  identity.ADT
	ActNr   seqnum.ADT
	CompRef compsem.SemRef
	ChnlID  identity.ADT
	ContExp termexp.ExpRec
}

func (r SubRec) turn() {}
func (r SubRec) GetConnID() identity.ADT { return r.ConnID }
func (r SubRec) GetActNr() seqnum.ADT { return r.ActNr }

func ErrRecTypeUnexpected(got TurnRec) error {
	return fmt.Errorf("step rec unexpected: %T%+v", got, got)
}

func ErrRecTypeMismatch(got, want TurnRec) error {
	return fmt.Errorf("step rec mismatch: want %T, got %T", want, got)
}

func ErrStepKindUnexpected(got turnKind) error {
	return fmt.Errorf("step kind unexpected: %v", got)
}
