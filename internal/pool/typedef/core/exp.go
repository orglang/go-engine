package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/polarity"
	"orglang/go-engine/internal/adt/uniqsym"
)

type ExpSpec interface {
	spec()
}

type OneSpec struct{}

func (OneSpec) spec() {}

type LinkSpec struct {
	TypeQN uniqsym.ADT
}

func (LinkSpec) spec() {}

// aka Internal Choice
type PlusSpec struct {
	ProcQNs []uniqsym.ADT
	ContExp ExpSpec
}

func (PlusSpec) spec() {}

// aka External Choice
type WithSpec struct {
	ProcQNs []uniqsym.ADT
	ContExp ExpSpec
}

func (WithSpec) spec() {}

type UpSpec struct {
	ContExp ExpSpec
}

func (UpSpec) spec() {}

type DownSpec struct {
	ContExp ExpSpec
}

func (DownSpec) spec() {}

type ExpRef interface {
	contkey.ContentKeyable
}

type OneRef struct {
	ExpCK contkey.ADT
}

func (r OneRef) Key() contkey.ADT { return r.ExpCK }

type LinkRef struct {
	ExpCK contkey.ADT
}

func (r LinkRef) Key() contkey.ADT { return r.ExpCK }

type PlusRef struct {
	ExpCK contkey.ADT
}

func (r PlusRef) Key() contkey.ADT { return r.ExpCK }

type WithRef struct {
	ExpCK contkey.ADT
}

func (r WithRef) Key() contkey.ADT { return r.ExpCK }

// aka Stype
type ExpRec interface {
	polarity.Polarizable
	contkey.ContentKeyable
}

type ProdRec interface {
	Next() contkey.ADT
}

type SumRec interface {
	Next(uniqsym.ADT) contkey.ADT
}

type OneRec struct {
	ExpCK contkey.ADT
}

func (OneRec) spec() {}

func (r OneRec) Key() contkey.ADT { return r.ExpCK }

func (OneRec) Pol() polarity.ADT { return polarity.Pos }

// aka TpName
type LinkRec struct {
	ExpCK  contkey.ADT
	TypeQN uniqsym.ADT
}

func (LinkRec) spec() {}

func (r LinkRec) Key() contkey.ADT { return r.ExpCK }

func (LinkRec) Pol() polarity.ADT { return polarity.Zero }

// aka Internal Choice
type PlusRec struct {
	ExpCK   contkey.ADT
	ProcQNs []uniqsym.ADT
	ContExp ExpRec
}

func (PlusRec) spec() {}

func (r PlusRec) Key() contkey.ADT { return r.ExpCK }

func (r PlusRec) Next() contkey.ADT { return r.ContExp.Key() }

func (PlusRec) Pol() polarity.ADT { return polarity.Pos }

// aka External Choice
type WithRec struct {
	ExpCK   contkey.ADT
	ProcQNs []uniqsym.ADT
	ContExp ExpRec
}

func (WithRec) spec() {}

func (r WithRec) Key() contkey.ADT { return r.ExpCK }

func (r WithRec) Next() contkey.ADT { return r.ContExp.Key() }

func (WithRec) Pol() polarity.ADT { return polarity.Neg }

type UpRec struct {
	ExpCK   contkey.ADT
	ContExp ExpRec
}

func (UpRec) spec() {}

func (r UpRec) Key() contkey.ADT { return r.ExpCK }

func (r UpRec) Next() contkey.ADT { return r.ContExp.Key() }

func (UpRec) Pol() polarity.ADT { return polarity.Zero }

type DownRec struct {
	ExpCK   contkey.ADT
	ContExp ExpRec
}

func (DownRec) spec() {}

func (r DownRec) Key() contkey.ADT { return r.ExpCK }

func (r DownRec) Next() contkey.ADT { return r.ContExp.Key() }

func (DownRec) Pol() polarity.ADT { return polarity.Zero }

func ErrSpecTypeUnexpected(got ExpSpec) error {
	return fmt.Errorf("spec type unexpected: %T", got)
}

func ErrRefTypeUnexpected(got ExpRef) error {
	return fmt.Errorf("ref type unexpected: %T", got)
}

func ErrRecTypeUnexpected(got ExpRec) error {
	return fmt.Errorf("rec type unexpected: %T", got)
}

func ErrDoesNotExist(want contkey.ADT) error {
	return fmt.Errorf("root doesn't exist: %v", want)
}
