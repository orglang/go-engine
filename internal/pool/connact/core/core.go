package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

// ExpRec is the expression-record contract carried by a connection action.
type ExpRec interface {
	IsExpRec()
}

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
	CompRef comproot.SemRef
	ChnlID  identity.ADT
	ValExp  ExpRec
}

func (r PubRec) act()                    {}
func (r PubRec) GetConnID() identity.ADT { return r.ConnID }
func (r PubRec) GetActNr() seqnum.ADT    { return r.ActNr }

// subscription (aka srv)
type SubRec struct {
	ConnID  identity.ADT
	ActNr   seqnum.ADT
	CompRef comproot.SemRef
	ChnlID  identity.ADT
	ContExp ExpRec
}

func (r SubRec) act()                    {}
func (r SubRec) GetConnID() identity.ADT { return r.ConnID }
func (r SubRec) GetActNr() seqnum.ADT    { return r.ActNr }

func ErrActRecTypeUnexpected(got ActRec) error {
	return fmt.Errorf("act rec unexpected: %T%+v", got, got)
}

func ErrActRecTypeMismatch(got, want ActRec) error {
	return fmt.Errorf("act rec mismatch: want %T, got %T", want, got)
}
