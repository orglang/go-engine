package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

type ExprSpec interface {
	Via() symbol.ADT
}

type CloseSpec struct {
	ContChnlPH symbol.ADT
}

func (s CloseSpec) Via() symbol.ADT { return s.ContChnlPH }

type WaitSpec struct {
	ContChnlPH symbol.ADT
	ContExpr    ExprSpec
}

func (s WaitSpec) Via() symbol.ADT { return s.ContChnlPH }

type SendSpec struct {
	CommChnlPH symbol.ADT
	ValChnlPH  symbol.ADT
}

func (s SendSpec) Via() symbol.ADT { return s.CommChnlPH }

type RecvSpec struct {
	CommChnlPH symbol.ADT
	NewChnlPH  symbol.ADT
	ContExpr    ExprSpec
}

func (s RecvSpec) Via() symbol.ADT { return s.CommChnlPH }

type LabSpec struct {
	CommChnlPH symbol.ADT
	ValLabQN   uniqsym.ADT
	ContExpr    ExprSpec
}

func (s LabSpec) Via() symbol.ADT { return s.CommChnlPH }

type CaseSpec struct {
	CommChnlPH symbol.ADT
	ContExprs  map[uniqsym.ADT]ExprSpec
}

func (s CaseSpec) Via() symbol.ADT { return s.CommChnlPH }

// aka ExprName
type LinkSpec struct {
	ProcTermQN uniqsym.ADT
	X          identity.ADT
	Ys         []identity.ADT
}

func (s LinkSpec) Via() symbol.ADT { return "" }

type FwdSpec struct {
	CommChnlPH symbol.ADT // old via | from | x
	ContChnlPH symbol.ADT // new via | to | y
}

func (s FwdSpec) Via() symbol.ADT { return s.CommChnlPH }

type CallSpec struct {
	CommChnlPH symbol.ADT // TODO удалить
	NewChnlPH  symbol.ADT
	ProcTermQN uniqsym.ADT
	ValChnlPHs []symbol.ADT // channel bulk
	ContExpr    ExprSpec
}

func (s CallSpec) Via() symbol.ADT { return s.CommChnlPH }

type SpawnSpec struct {
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	NewChnlPHs []symbol.ADT
	ContExpr    ExprSpec
}

func (s SpawnSpec) Via() symbol.ADT { return s.CommChnlPH }

type AcqureSpec struct {
	CommChnlPH symbol.ADT
	ContExpr    ExprSpec
}

func (s AcqureSpec) Via() symbol.ADT { return s.CommChnlPH }

type AcceptSpec struct {
	CommChnlPH symbol.ADT
	ContExpr    ExprSpec
}

func (s AcceptSpec) Via() symbol.ADT { return s.CommChnlPH }

type DetachSpec struct {
	CommChnlPH symbol.ADT
}

func (s DetachSpec) Via() symbol.ADT { return s.CommChnlPH }

type ReleaseSpec struct {
	CommChnlPH symbol.ADT
}

func (s ReleaseSpec) Via() symbol.ADT { return s.CommChnlPH }

type ExprRec interface {
	ExprSpec
	IsExprRec()
}

type CloseRec struct {
	ContChnlPH symbol.ADT
}

func (r CloseRec) Via() symbol.ADT { return r.ContChnlPH }

func (CloseRec) IsExprRec() {}

type WaitRec struct {
	ContChnlPH symbol.ADT
	ContExpr    ExprSpec
}

func (r WaitRec) Via() symbol.ADT { return r.ContChnlPH }

func (WaitRec) IsExprRec() {}

type SendRec struct {
	CommChnlPH symbol.ADT
	ContChnlID identity.ADT
	ValChnlID  identity.ADT
	ValExprCK  contkey.ADT
	CommID     identity.ADT
}

func (r SendRec) Via() symbol.ADT { return r.CommChnlPH }

func (SendRec) IsExprRec() {}

type RecvRec struct {
	CommChnlPH symbol.ADT
	ContChnlID identity.ADT
	NewChnlPH  symbol.ADT
	ContExpr    ExprSpec
}

func (r RecvRec) Via() symbol.ADT { return r.CommChnlPH }

func (RecvRec) IsExprRec() {}

type LabRec struct {
	CommChnlPH symbol.ADT
	ContChnlID identity.ADT
	ValExprCK  contkey.ADT
	ValLabQN   uniqsym.ADT // TODO выпилить
}

func (r LabRec) Via() symbol.ADT { return r.CommChnlPH }

func (LabRec) IsExprRec() {}

type CaseRec struct {
	CommChnlPH symbol.ADT
	ContChnlID identity.ADT
	ContExprs  map[uniqsym.ADT]ExprSpec
}

func (r CaseRec) Via() symbol.ADT { return r.CommChnlPH }

func (CaseRec) IsExprRec() {}

type FwdRec struct {
	CommChnlPH symbol.ADT
	ContChnlID identity.ADT
}

func (r FwdRec) Via() symbol.ADT { return r.CommChnlPH }

func (FwdRec) IsExprRec() {}

func CollectEnv(spec ExprSpec) []identity.ADT {
	return collectEnvRec(spec, []identity.ADT{})
}

func collectEnvRec(s ExprSpec, env []identity.ADT) []identity.ADT {
	switch spec := s.(type) {
	case RecvSpec:
		return collectEnvRec(spec.ContExpr, env)
	case CaseSpec:
		for _, cont := range spec.ContExprs {
			env = collectEnvRec(cont, env)
		}
		return env
	default:
		return env
	}
}

func ErrExprTypeUnexpected(got ExprSpec) error {
	return fmt.Errorf("expr spec unexpected: %T", got)
}

func ErrRecTypeUnexpected(got any) error {
	return fmt.Errorf("expr rec unexpected: %T", got)
}

func ErrExprTypeMismatch(got, want ExprSpec) error {
	return fmt.Errorf("expr spec mismatch: want %T, got %T", want, got)
}

func ErrExprValueNil(pid identity.ADT) error {
	return fmt.Errorf("proc %q term is nil", pid)
}
