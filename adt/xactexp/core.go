package xactexp

import (
	"fmt"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/polarity"
	"orglang/go-engine/adt/uniqsym"
)

type ExpSpec interface {
	spec()
}

type OneSpec struct{}

func (OneSpec) spec() {}

type LinkSpec struct {
	XactQN uniqsym.ADT
}

func (LinkSpec) spec() {}

// aka Internal Choice
type PlusSpec struct {
	Choices map[uniqsym.ADT]ExpSpec // conts
}

func (PlusSpec) spec() {}

// aka External Choice
type WithSpec struct {
	Choices map[uniqsym.ADT]ExpSpec // conts
}

func (WithSpec) spec() {}

type ExpRef interface {
	identity.Identifiable
}

type OneRef struct {
	ExpID identity.ADT
}

func (r OneRef) Ident() identity.ADT { return r.ExpID }

type LinkRef struct {
	ExpID identity.ADT
}

func (r LinkRef) Ident() identity.ADT { return r.ExpID }

type PlusRef struct {
	ExpID identity.ADT
}

func (r PlusRef) Ident() identity.ADT { return r.ExpID }

type WithRef struct {
	ExpID identity.ADT
}

func (r WithRef) Ident() identity.ADT { return r.ExpID }

// aka Stype
type ExpRec interface {
	identity.Identifiable
	polarity.Polarizable
}

type ProdRec interface {
	Next() identity.ADT
}

type SumRec interface {
	Next(uniqsym.ADT) identity.ADT
}

type OneRec struct {
	ExpID identity.ADT
}

func (OneRec) spec() {}

func (r OneRec) Ident() identity.ADT { return r.ExpID }

func (OneRec) Pol() polarity.ADT { return polarity.Pos }

// aka TpName
type LinkRec struct {
	ExpID  identity.ADT
	XactQN uniqsym.ADT
}

func (LinkRec) spec() {}

func (r LinkRec) Ident() identity.ADT { return r.ExpID }

func (LinkRec) Pol() polarity.ADT { return polarity.Zero }

// aka Internal Choice
type PlusRec struct {
	ExpID identity.ADT
	Zs    map[uniqsym.ADT]ExpRec
}

func (PlusRec) spec() {}

func (r PlusRec) Ident() identity.ADT { return r.ExpID }

func (r PlusRec) Next(l uniqsym.ADT) identity.ADT { return r.Zs[l].Ident() }

func (PlusRec) Pol() polarity.ADT { return polarity.Pos }

// aka External Choice
type WithRec struct {
	ExpID identity.ADT
	Zs    map[uniqsym.ADT]ExpRec
}

func (WithRec) spec() {}

func (r WithRec) Ident() identity.ADT { return r.ExpID }

func (r WithRec) Next(l uniqsym.ADT) identity.ADT { return r.Zs[l].Ident() }

func (WithRec) Pol() polarity.ADT { return polarity.Neg }

func ErrSpecTypeUnexpected(got ExpSpec) error {
	return fmt.Errorf("spec type unexpected: %T", got)
}

func ErrRefTypeUnexpected(got ExpRef) error {
	return fmt.Errorf("ref type unexpected: %T", got)
}

func ErrRecTypeUnexpected(got ExpRec) error {
	return fmt.Errorf("rec type unexpected: %T", got)
}
