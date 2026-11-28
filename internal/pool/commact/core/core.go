package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
)

// ExpRec is the expression-record contract carried by a communication action.
type ExpRec interface {
	IsExpRec()
}

// step (aka Act)
type ActRec interface {
	act()
	GetCommID() identity.ADT
	GetActNr() seqnum.ADT
}

// publication (aka msg)
type PubRec struct {
	CommID  identity.ADT
	ActNr   seqnum.ADT
	RootRef execroot.RootRef
	ChnlID  identity.ADT
	ValExp  ExpRec
}

func (r PubRec) act()                    {}
func (r PubRec) GetCommID() identity.ADT { return r.CommID }
func (r PubRec) GetActNr() seqnum.ADT    { return r.ActNr }

// subscription (aka srv)
type SubRec struct {
	CommID  identity.ADT
	ActNr   seqnum.ADT
	RootRef execroot.RootRef
	ChnlID  identity.ADT
	ContExp ExpRec
}

func (r SubRec) act()                    {}
func (r SubRec) GetCommID() identity.ADT { return r.CommID }
func (r SubRec) GetActNr() seqnum.ADT    { return r.ActNr }

func ErrActRecTypeUnexpected(got ActRec) error {
	return fmt.Errorf("act rec unexpected: %T%+v", got, got)
}

func ErrActRecTypeMismatch(got, want ActRec) error {
	return fmt.Errorf("act rec mismatch: want %T, got %T", want, got)
}
