package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

// ExprSpec is one node in a Term Definition tree.
type ExprSpec interface{ exprSpec() }
type ExprAcquireSpec struct {
	CommChnlPH symbol.ADT
	ContExp    ExprSpec
}

func (ExprAcquireSpec) exprSpec() {}

type ExprAcceptSpec struct {
	CommChnlPH symbol.ADT
	ContExp    ExprSpec
}

func (ExprAcceptSpec) exprSpec() {}

type ExprHireSpec struct {
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExprSpec
}

func (ExprHireSpec) exprSpec() {}

type ExprApplySpec struct {
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExprSpec
}

func (ExprApplySpec) exprSpec() {}

type ExprReleaseSpec struct{ CommChnlPH symbol.ADT }

func (ExprReleaseSpec) exprSpec() {}

type ExprDetachSpec struct{ CommChnlPH symbol.ADT }

func (ExprDetachSpec) exprSpec() {}

type ExprFireSpec struct {
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (ExprFireSpec) exprSpec() {}

type ExprQuitSpec struct {
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (ExprQuitSpec) exprSpec() {}

type ExprSpawnSpec struct {
	ProcTermQN  uniqsym.ADT
	ProcCompQNs []uniqsym.ADT
}

func (ExprSpawnSpec) exprSpec() {}

// ExprRec is the recursive, content-addressed domain representation of a Definition.
type ExprRec interface {
	exprRec()
	ContentKey() contkey.ADT
}
type ExprAcquireRec struct {
	ExprCK      contkey.ADT
	CommChnlPH symbol.ADT
	ContExp    ExprRec
}

func (ExprAcquireRec) exprRec()                   {}
func (r ExprAcquireRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprAcceptRec struct {
	ExprCK      contkey.ADT
	CommChnlPH symbol.ADT
	ContExp    ExprRec
}

func (ExprAcceptRec) exprRec()                   {}
func (r ExprAcceptRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprHireRec struct {
	ExprCK      contkey.ADT
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExprRec
}

func (ExprHireRec) exprRec()                   {}
func (r ExprHireRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprApplyRec struct {
	ExprCK      contkey.ADT
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExprRec
}

func (ExprApplyRec) exprRec()                   {}
func (r ExprApplyRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprReleaseRec struct {
	ExprCK      contkey.ADT
	CommChnlPH symbol.ADT
}

func (ExprReleaseRec) exprRec()                   {}
func (r ExprReleaseRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprDetachRec struct {
	ExprCK      contkey.ADT
	CommChnlPH symbol.ADT
}

func (ExprDetachRec) exprRec()                   {}
func (r ExprDetachRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprFireRec struct {
	ExprCK      contkey.ADT
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (ExprFireRec) exprRec()                   {}
func (r ExprFireRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprQuitRec struct {
	ExprCK      contkey.ADT
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (ExprQuitRec) exprRec()                   {}
func (r ExprQuitRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprSpawnRec struct {
	ExprCK       contkey.ADT
	ProcTermQN  uniqsym.ADT
	ProcCompQNs []uniqsym.ADT
}

func (ExprSpawnRec) exprRec()                   {}
func (r ExprSpawnRec) ContentKey() contkey.ADT { return r.ExprCK }

func ConvertSpecToRec(spec ExprSpec) (ExprRec, error) {
	key, _, err := ExprContentKey(spec)
	if err != nil {
		return nil, err
	}
	switch s := spec.(type) {
	case ExprAcquireSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExprAcquireRec{key, s.CommChnlPH, c}, nil
	case ExprAcceptSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExprAcceptRec{key, s.CommChnlPH, c}, nil
	case ExprHireSpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExprHireRec{key, s.CommChnlPH, s.ProcTermQN, c}, nil
	case ExprApplySpec:
		c, e := ConvertSpecToRec(s.ContExp)
		if e != nil {
			return nil, e
		}
		return ExprApplyRec{key, s.CommChnlPH, s.ProcTermQN, c}, nil
	case ExprReleaseSpec:
		return ExprReleaseRec{key, s.CommChnlPH}, nil
	case ExprDetachSpec:
		return ExprDetachRec{key, s.CommChnlPH}, nil
	case ExprFireSpec:
		return ExprFireRec{key, s.CommChnlPH, s.ProcTermQN}, nil
	case ExprQuitSpec:
		return ExprQuitRec{key, s.CommChnlPH, s.ProcTermQN}, nil
	case ExprSpawnSpec:
		return ExprSpawnRec{key, s.ProcTermQN, s.ProcCompQNs}, nil
	default:
		return nil, fmt.Errorf("unexpected pool Definition expression spec %T", spec)
	}
}
func ConvertRecToSpec(rec ExprRec) ExprSpec {
	switch r := rec.(type) {
	case ExprAcquireRec:
		return ExprAcquireSpec{r.CommChnlPH, ConvertRecToSpec(r.ContExp)}
	case ExprAcceptRec:
		return ExprAcceptSpec{r.CommChnlPH, ConvertRecToSpec(r.ContExp)}
	case ExprHireRec:
		return ExprHireSpec{r.CommChnlPH, r.ProcTermQN, ConvertRecToSpec(r.ContExp)}
	case ExprApplyRec:
		return ExprApplySpec{r.CommChnlPH, r.ProcTermQN, ConvertRecToSpec(r.ContExp)}
	case ExprReleaseRec:
		return ExprReleaseSpec{r.CommChnlPH}
	case ExprDetachRec:
		return ExprDetachSpec{r.CommChnlPH}
	case ExprFireRec:
		return ExprFireSpec{r.CommChnlPH, r.ProcTermQN}
	case ExprQuitRec:
		return ExprQuitSpec{r.CommChnlPH, r.ProcTermQN}
	case ExprSpawnRec:
		return ExprSpawnSpec{r.ProcTermQN, r.ProcCompQNs}
	default:
		panic(fmt.Sprintf("unexpected pool Definition expression rec %T", rec))
	}
}
