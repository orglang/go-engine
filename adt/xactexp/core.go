package xactexp

import (
	"fmt"

	"orglang/go-engine/adt/polarity"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
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
	valkey.Keyable
}

type OneRef struct {
	ExpVK valkey.ADT
}

func (r OneRef) Key() valkey.ADT { return r.ExpVK }

type LinkRef struct {
	ExpVK valkey.ADT
}

func (r LinkRef) Key() valkey.ADT { return r.ExpVK }

type PlusRef struct {
	ExpVK valkey.ADT
}

func (r PlusRef) Key() valkey.ADT { return r.ExpVK }

type WithRef struct {
	ExpVK valkey.ADT
}

func (r WithRef) Key() valkey.ADT { return r.ExpVK }

// aka Stype
type ExpRec interface {
	polarity.Polarizable
	valkey.Keyable
}

type ProdRec interface {
	Next() valkey.ADT
}

type SumRec interface {
	Next(uniqsym.ADT) valkey.ADT
}

type OneRec struct {
	ExpVK valkey.ADT
}

func (OneRec) spec() {}

func (r OneRec) Key() valkey.ADT { return r.ExpVK }

func (OneRec) Pol() polarity.ADT { return polarity.Pos }

// aka TpName
type LinkRec struct {
	ExpVK  valkey.ADT
	XactQN uniqsym.ADT
}

func (LinkRec) spec() {}

func (r LinkRec) Key() valkey.ADT { return r.ExpVK }

func (LinkRec) Pol() polarity.ADT { return polarity.Zero }

// aka Internal Choice
type PlusRec struct {
	ExpVK   valkey.ADT
	Choices map[uniqsym.ADT]ExpRec
}

func (PlusRec) spec() {}

func (r PlusRec) Key() valkey.ADT { return r.ExpVK }

func (r PlusRec) Next(l uniqsym.ADT) valkey.ADT { return r.Choices[l].Key() }

func (PlusRec) Pol() polarity.ADT { return polarity.Pos }

// aka External Choice
type WithRec struct {
	ExpVK   valkey.ADT
	Choices map[uniqsym.ADT]ExpRec
}

func (WithRec) spec() {}

func (r WithRec) Key() valkey.ADT { return r.ExpVK }

func (r WithRec) Next(l uniqsym.ADT) valkey.ADT { return r.Choices[l].Key() }

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
