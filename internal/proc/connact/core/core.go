package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
)

// ExpSpec is the process expression contract used by execution and connection actions.
type ExpSpec interface {
	Via() symbol.ADT
}

// ExpRec is a persisted process expression record carried by a connection action.
type ExpRec interface {
	ExpSpec
	IsExpRec()
}

// aka Sem
type ActRec interface {
	act() identity.ADT
	GetConnID() identity.ADT
	GetActNr() seqnum.ADT
}

func ChnlID(rec ActRec) identity.ADT { return rec.act() }

type PubRec struct {
	ConnID  identity.ADT
	ActNr   seqnum.ADT
	CompRef comproot.SemRef
	ChnlID  identity.ADT
	ValExp  ExpRec
}

func (r PubRec) act() identity.ADT      { return r.ChnlID }
func (r PubRec) GetConnID() identity.ADT { return r.ConnID }
func (r PubRec) GetActNr() seqnum.ADT    { return r.ActNr }

type SubRec struct {
	ConnID  identity.ADT
	ActNr   seqnum.ADT
	CompRef comproot.SemRef
	ChnlID  identity.ADT
	ContExp ExpRec
}

func (r SubRec) act() identity.ADT      { return r.ChnlID }
func (r SubRec) GetConnID() identity.ADT { return r.ConnID }
func (r SubRec) GetActNr() seqnum.ADT    { return r.ActNr }

func ErrActRecTypeUnexpected(got ActRec) error {
	return fmt.Errorf("act rec unexpected: %T", got)
}

func ErrActRecTypeMismatch(got, want ActRec) error {
	return fmt.Errorf("act rec mismatch: want %T, got %T", want, got)
}
