package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

// ExprSpec is one node in a Term Definition tree.
type ExprSpec interface {
	exprSpec()
}

type ExprCloseSpec struct{ ContChnlPH symbol.ADT }

func (ExprCloseSpec) exprSpec() {}

type ExprWaitSpec struct {
	ContChnlPH symbol.ADT
	ContExpr   ExprSpec
}

func (ExprWaitSpec) exprSpec() {}

type ExprSendSpec struct {
	CommChnlPH symbol.ADT
	ValChnlPH  symbol.ADT
}

func (ExprSendSpec) exprSpec() {}

type ExprRecvSpec struct {
	CommChnlPH symbol.ADT
	NewChnlPH  symbol.ADT
	ContExpr   ExprSpec
}

func (ExprRecvSpec) exprSpec() {}

type ExprLabSpec struct {
	CommChnlPH symbol.ADT
	ValLabQN   uniqsym.ADT
	ContExpr   ExprSpec
}

func (ExprLabSpec) exprSpec() {}

type ExprCaseSpec struct {
	CommChnlPH symbol.ADT
	ContExprs  map[uniqsym.ADT]ExprSpec
}

func (ExprCaseSpec) exprSpec() {}

type ExprLinkSpec struct {
	ProcTermQN uniqsym.ADT
	X          identity.ADT
	Ys         []identity.ADT
}

func (ExprLinkSpec) exprSpec() {}

type ExprFwdSpec struct {
	CommChnlPH symbol.ADT
	ContChnlPH symbol.ADT
}

func (ExprFwdSpec) exprSpec() {}

type ExprCallSpec struct {
	CommChnlPH symbol.ADT
	NewChnlPH  symbol.ADT
	ProcTermQN uniqsym.ADT
	ValChnlPHs []symbol.ADT
	ContExpr   ExprSpec
}

func (ExprCallSpec) exprSpec() {}

type ExprSpawnSpec struct {
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	NewChnlPHs []symbol.ADT
	ContExpr   ExprSpec
}

func (ExprSpawnSpec) exprSpec() {}

type ExprAcquireSpec struct {
	CommChnlPH symbol.ADT
	ContExpr   ExprSpec
}

func (ExprAcquireSpec) exprSpec() {}

type ExprAcceptSpec struct {
	CommChnlPH symbol.ADT
	ContExpr   ExprSpec
}

func (ExprAcceptSpec) exprSpec() {}

type ExprDetachSpec struct{ CommChnlPH symbol.ADT }

func (ExprDetachSpec) exprSpec() {}

type ExprReleaseSpec struct{ CommChnlPH symbol.ADT }

func (ExprReleaseSpec) exprSpec() {}

// ExprRec is the recursive, content-addressed domain representation of a Definition.
type ExprRec interface {
	exprRec()
	ContentKey() contkey.ADT
}
type ExprCloseRec struct {
	ExprCK     contkey.ADT
	ContChnlPH symbol.ADT
}

func (ExprCloseRec) exprRec()                  {}
func (r ExprCloseRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprWaitRec struct {
	ExprCK     contkey.ADT
	ContChnlPH symbol.ADT
	ContExpr   ExprRec
}

func (ExprWaitRec) exprRec()                  {}
func (r ExprWaitRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprSendRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
	ValChnlPH  symbol.ADT
}

func (ExprSendRec) exprRec()                  {}
func (r ExprSendRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprRecvRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
	NewChnlPH  symbol.ADT
	ContExpr   ExprRec
}

func (ExprRecvRec) exprRec()                  {}
func (r ExprRecvRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprLabRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
	ValLabQN   uniqsym.ADT
	ContExpr   ExprRec
}

func (ExprLabRec) exprRec()                  {}
func (r ExprLabRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprCaseRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
	ContExprs  map[uniqsym.ADT]ExprRec
}

func (ExprCaseRec) exprRec()                  {}
func (r ExprCaseRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprLinkRec struct {
	ExprCK     contkey.ADT
	ProcTermQN uniqsym.ADT
	X          identity.ADT
	Ys         []identity.ADT
}

func (ExprLinkRec) exprRec()                  {}
func (r ExprLinkRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprFwdRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
	ContChnlPH symbol.ADT
}

func (ExprFwdRec) exprRec()                  {}
func (r ExprFwdRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprCallRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
	NewChnlPH  symbol.ADT
	ProcTermQN uniqsym.ADT
	ValChnlPHs []symbol.ADT
	ContExpr   ExprRec
}

func (ExprCallRec) exprRec()                  {}
func (r ExprCallRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprSpawnRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	NewChnlPHs []symbol.ADT
	ContExpr   ExprRec
}

func (ExprSpawnRec) exprRec()                  {}
func (r ExprSpawnRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprAcquireRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
	ContExpr   ExprRec
}

func (ExprAcquireRec) exprRec()                  {}
func (r ExprAcquireRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprAcceptRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
	ContExpr   ExprRec
}

func (ExprAcceptRec) exprRec()                  {}
func (r ExprAcceptRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprDetachRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
}

func (ExprDetachRec) exprRec()                  {}
func (r ExprDetachRec) ContentKey() contkey.ADT { return r.ExprCK }

type ExprReleaseRec struct {
	ExprCK     contkey.ADT
	CommChnlPH symbol.ADT
}

func (ExprReleaseRec) exprRec()                  {}
func (r ExprReleaseRec) ContentKey() contkey.ADT { return r.ExprCK }

// ConvertSpecToRec converts a declarative Definition into its content-addressed domain record.
func ConvertSpecToRec(spec ExprSpec) (ExprRec, error) {
	key, _, err := ExprContentKey(spec)
	if err != nil {
		return nil, err
	}
	switch s := spec.(type) {
	case ExprCloseSpec:
		return ExprCloseRec{key, s.ContChnlPH}, nil
	case ExprWaitSpec:
		c, e := ConvertSpecToRec(s.ContExpr)
		if e != nil {
			return nil, e
		}
		return ExprWaitRec{key, s.ContChnlPH, c}, nil
	case ExprSendSpec:
		return ExprSendRec{key, s.CommChnlPH, s.ValChnlPH}, nil
	case ExprRecvSpec:
		c, e := ConvertSpecToRec(s.ContExpr)
		if e != nil {
			return nil, e
		}
		return ExprRecvRec{key, s.CommChnlPH, s.NewChnlPH, c}, nil
	case ExprLabSpec:
		c, e := ConvertSpecToRec(s.ContExpr)
		if e != nil {
			return nil, e
		}
		return ExprLabRec{key, s.CommChnlPH, s.ValLabQN, c}, nil
	case ExprCaseSpec:
		cs := make(map[uniqsym.ADT]ExprRec, len(s.ContExprs))
		for q, v := range s.ContExprs {
			c, e := ConvertSpecToRec(v)
			if e != nil {
				return nil, e
			}
			cs[q] = c
		}
		return ExprCaseRec{key, s.CommChnlPH, cs}, nil
	case ExprLinkSpec:
		return ExprLinkRec{key, s.ProcTermQN, s.X, s.Ys}, nil
	case ExprFwdSpec:
		return ExprFwdRec{key, s.CommChnlPH, s.ContChnlPH}, nil
	case ExprCallSpec:
		c, e := ConvertSpecToRec(s.ContExpr)
		if e != nil {
			return nil, e
		}
		return ExprCallRec{key, s.CommChnlPH, s.NewChnlPH, s.ProcTermQN, s.ValChnlPHs, c}, nil
	case ExprSpawnSpec:
		c, e := ConvertSpecToRec(s.ContExpr)
		if e != nil {
			return nil, e
		}
		return ExprSpawnRec{key, s.CommChnlPH, s.ProcTermQN, s.NewChnlPHs, c}, nil
	case ExprAcquireSpec:
		c, e := ConvertSpecToRec(s.ContExpr)
		if e != nil {
			return nil, e
		}
		return ExprAcquireRec{key, s.CommChnlPH, c}, nil
	case ExprAcceptSpec:
		c, e := ConvertSpecToRec(s.ContExpr)
		if e != nil {
			return nil, e
		}
		return ExprAcceptRec{key, s.CommChnlPH, c}, nil
	case ExprDetachSpec:
		return ExprDetachRec{key, s.CommChnlPH}, nil
	case ExprReleaseSpec:
		return ExprReleaseRec{key, s.CommChnlPH}, nil
	default:
		return nil, fmt.Errorf("unexpected process Definition expression spec %T", spec)
	}
}

func ConvertRecToSpec(rec ExprRec) ExprSpec {
	switch r := rec.(type) {
	case ExprCloseRec:
		return ExprCloseSpec{r.ContChnlPH}
	case ExprWaitRec:
		return ExprWaitSpec{r.ContChnlPH, ConvertRecToSpec(r.ContExpr)}
	case ExprSendRec:
		return ExprSendSpec{r.CommChnlPH, r.ValChnlPH}
	case ExprRecvRec:
		return ExprRecvSpec{r.CommChnlPH, r.NewChnlPH, ConvertRecToSpec(r.ContExpr)}
	case ExprLabRec:
		return ExprLabSpec{r.CommChnlPH, r.ValLabQN, ConvertRecToSpec(r.ContExpr)}
	case ExprCaseRec:
		cs := make(map[uniqsym.ADT]ExprSpec, len(r.ContExprs))
		for q, v := range r.ContExprs {
			cs[q] = ConvertRecToSpec(v)
		}
		return ExprCaseSpec{r.CommChnlPH, cs}
	case ExprLinkRec:
		return ExprLinkSpec{r.ProcTermQN, r.X, r.Ys}
	case ExprFwdRec:
		return ExprFwdSpec{r.CommChnlPH, r.ContChnlPH}
	case ExprCallRec:
		return ExprCallSpec{r.CommChnlPH, r.NewChnlPH, r.ProcTermQN, r.ValChnlPHs, ConvertRecToSpec(r.ContExpr)}
	case ExprSpawnRec:
		return ExprSpawnSpec{r.CommChnlPH, r.ProcTermQN, r.NewChnlPHs, ConvertRecToSpec(r.ContExpr)}
	case ExprAcquireRec:
		return ExprAcquireSpec{r.CommChnlPH, ConvertRecToSpec(r.ContExpr)}
	case ExprAcceptRec:
		return ExprAcceptSpec{r.CommChnlPH, ConvertRecToSpec(r.ContExpr)}
	case ExprDetachRec:
		return ExprDetachSpec{r.CommChnlPH}
	case ExprReleaseRec:
		return ExprReleaseSpec{r.CommChnlPH}
	default:
		panic(fmt.Sprintf("unexpected process Definition expression rec %T", rec))
	}
}
