package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

// ExpSpec is one node in a Term Definition tree.
type ExpSpec interface{ expSpec() }
type ExpAcquireSpec struct {
	ConnChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (ExpAcquireSpec) expSpec() {}

type ExpAcceptSpec struct {
	ConnChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (ExpAcceptSpec) expSpec() {}

type ExpHireSpec struct {
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExpSpec
}

func (ExpHireSpec) expSpec() {}

type ExpApplySpec struct {
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExpSpec
}

func (ExpApplySpec) expSpec() {}

type ExpReleaseSpec struct{ ConnChnlPH symbol.ADT }

func (ExpReleaseSpec) expSpec() {}

type ExpDetachSpec struct{ ConnChnlPH symbol.ADT }

func (ExpDetachSpec) expSpec() {}

type ExpFireSpec struct {
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (ExpFireSpec) expSpec() {}

type ExpQuitSpec struct {
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (ExpQuitSpec) expSpec() {}

type ExpSpawnSpec struct {
	ProcTermQN  uniqsym.ADT
	ProcCompQNs []uniqsym.ADT
}

func (ExpSpawnSpec) expSpec() {}

// ExpRec is the recursive, content-addressed domain representation of a Definition.
type ExpRec interface {
	expRec()
	ContentKey() contkey.ADT
}
type ExpAcquireRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ContExp    ExpRec
}

func (ExpAcquireRec) expRec()                   {}
func (r ExpAcquireRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpAcceptRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ContExp    ExpRec
}

func (ExpAcceptRec) expRec()                   {}
func (r ExpAcceptRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpHireRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExpRec
}

func (ExpHireRec) expRec()                   {}
func (r ExpHireRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpApplyRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExpRec
}

func (ExpApplyRec) expRec()                   {}
func (r ExpApplyRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpReleaseRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
}

func (ExpReleaseRec) expRec()                   {}
func (r ExpReleaseRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpDetachRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
}

func (ExpDetachRec) expRec()                   {}
func (r ExpDetachRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpFireRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (ExpFireRec) expRec()                   {}
func (r ExpFireRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpQuitRec struct {
	ExpCK      contkey.ADT
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (ExpQuitRec) expRec()                   {}
func (r ExpQuitRec) ContentKey() contkey.ADT { return r.ExpCK }

type ExpSpawnRec struct {
	ExpCK       contkey.ADT
	ProcTermQN  uniqsym.ADT
	ProcCompQNs []uniqsym.ADT
}

func (ExpSpawnRec) expRec()                   {}
func (r ExpSpawnRec) ContentKey() contkey.ADT { return r.ExpCK }

func ConvertSpecToRec(spec ExpSpec) (ExpRec, error) {
	key, _, err := ExpContentKey(spec)
	if err != nil {
		return nil, err
	}
	switch s := spec.(type) {
	case ExpAcquireSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpAcquireRec{key, s.ConnChnlPH, c}, nil
	case ExpAcceptSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpAcceptRec{key, s.ConnChnlPH, c}, nil
	case ExpHireSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpHireRec{key, s.ConnChnlPH, s.ProcTermQN, c}, nil
	case ExpApplySpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExpApplyRec{key, s.ConnChnlPH, s.ProcTermQN, c}, nil
	case ExpReleaseSpec:
		return ExpReleaseRec{key, s.ConnChnlPH}, nil
	case ExpDetachSpec:
		return ExpDetachRec{key, s.ConnChnlPH}, nil
	case ExpFireSpec:
		return ExpFireRec{key, s.ConnChnlPH, s.ProcTermQN}, nil
	case ExpQuitSpec:
		return ExpQuitRec{key, s.ConnChnlPH, s.ProcTermQN}, nil
	case ExpSpawnSpec:
		return ExpSpawnRec{key, s.ProcTermQN, s.ProcCompQNs}, nil
	default:
		return nil, fmt.Errorf("unexpected pool Definition expression spec %T", spec)
	}
}
func ConvertRecToSpec(rec ExpRec) ExpSpec {
	switch r := rec.(type) {
	case ExpAcquireRec:
		return ExpAcquireSpec{r.ConnChnlPH, ConvertRecToSpec(r.ContExp)}
	case ExpAcceptRec:
		return ExpAcceptSpec{r.ConnChnlPH, ConvertRecToSpec(r.ContExp)}
	case ExpHireRec:
		return ExpHireSpec{r.ConnChnlPH, r.ProcTermQN, ConvertRecToSpec(r.ContExp)}
	case ExpApplyRec:
		return ExpApplySpec{r.ConnChnlPH, r.ProcTermQN, ConvertRecToSpec(r.ContExp)}
	case ExpReleaseRec:
		return ExpReleaseSpec{r.ConnChnlPH}
	case ExpDetachRec:
		return ExpDetachSpec{r.ConnChnlPH}
	case ExpFireRec:
		return ExpFireSpec{r.ConnChnlPH, r.ProcTermQN}
	case ExpQuitRec:
		return ExpQuitSpec{r.ConnChnlPH, r.ProcTermQN}
	case ExpSpawnRec:
		return ExpSpawnSpec{r.ProcTermQN, r.ProcCompQNs}
	default:
		panic(fmt.Sprintf("unexpected pool Definition expression rec %T", rec))
	}
}
