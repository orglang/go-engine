package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

// step (aka Act)
type ActRec interface {
	act()
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

func (r PubRec) act()                   {}
func (r PubRec) GetConnID() identity.ADT { return r.ConnID }
func (r PubRec) GetActNr() seqnum.ADT    { return r.ActNr }

// subscription (aka srv)
type SubRec struct {
	ConnID  identity.ADT
	ActNr   seqnum.ADT
	CompRef compsem.SemRef
	ChnlID  identity.ADT
	ContExp termexp.ExpRec
}

func (r SubRec) act()                   {}
func (r SubRec) GetConnID() identity.ADT { return r.ConnID }
func (r SubRec) GetActNr() seqnum.ADT    { return r.ActNr }

func ErrActRecTypeUnexpected(got ActRec) error {
	return fmt.Errorf("act rec unexpected: %T%+v", got, got)
}

func ErrActRecTypeMismatch(got, want ActRec) error {
	return fmt.Errorf("act rec mismatch: want %T, got %T", want, got)
}

func ErrActKindUnexpected(got actKind) error {
	return fmt.Errorf("act kind unexpected: %v", got)
}
