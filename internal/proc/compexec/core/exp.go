package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

type ExpSpec interface {
	Via() symbol.ADT
}

type CloseSpec struct {
	ContChnlPH symbol.ADT
}

func (s CloseSpec) Via() symbol.ADT { return s.ContChnlPH }

type WaitSpec struct {
	ContChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (s WaitSpec) Via() symbol.ADT { return s.ContChnlPH }

type SendSpec struct {
	ConnChnlPH symbol.ADT
	ValChnlPH  symbol.ADT
}

func (s SendSpec) Via() symbol.ADT { return s.ConnChnlPH }

type RecvSpec struct {
	ConnChnlPH symbol.ADT
	NewChnlPH  symbol.ADT
	ContExp    ExpSpec
}

func (s RecvSpec) Via() symbol.ADT { return s.ConnChnlPH }

type LabSpec struct {
	ConnChnlPH symbol.ADT
	ValLabQN   uniqsym.ADT
	ContExp    ExpSpec
}

func (s LabSpec) Via() symbol.ADT { return s.ConnChnlPH }

type CaseSpec struct {
	ConnChnlPH symbol.ADT
	ContExps   map[uniqsym.ADT]ExpSpec
}

func (s CaseSpec) Via() symbol.ADT { return s.ConnChnlPH }

// aka ExpName
type LinkSpec struct {
	ProcTermQN uniqsym.ADT
	X          identity.ADT
	Ys         []identity.ADT
}

func (s LinkSpec) Via() symbol.ADT { return "" }

type FwdSpec struct {
	ConnChnlPH symbol.ADT // old via | from | x
	ContChnlPH symbol.ADT // new via | to | y
}

func (s FwdSpec) Via() symbol.ADT { return s.ConnChnlPH }

type CallSpec struct {
	ConnChnlPH symbol.ADT // TODO удалить
	NewChnlPH  symbol.ADT
	ProcTermQN uniqsym.ADT
	ValChnlPHs []symbol.ADT // channel bulk
	ContExp    ExpSpec
}

func (s CallSpec) Via() symbol.ADT { return s.ConnChnlPH }

type SpawnSpec struct {
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	NewChnlPHs []symbol.ADT
	ContExp    ExpSpec
}

func (s SpawnSpec) Via() symbol.ADT { return s.ConnChnlPH }

type AcqureSpec struct {
	ConnChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (s AcqureSpec) Via() symbol.ADT { return s.ConnChnlPH }

type AcceptSpec struct {
	ConnChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (s AcceptSpec) Via() symbol.ADT { return s.ConnChnlPH }

type DetachSpec struct {
	ConnChnlPH symbol.ADT
}

func (s DetachSpec) Via() symbol.ADT { return s.ConnChnlPH }

type ReleaseSpec struct {
	ConnChnlPH symbol.ADT
}

func (s ReleaseSpec) Via() symbol.ADT { return s.ConnChnlPH }

type ExpRec interface {
	ExpSpec
	impl()
}

type CloseRec struct {
	ContChnlPH symbol.ADT
}

func (r CloseRec) Via() symbol.ADT { return r.ContChnlPH }

func (CloseRec) impl() {}

type WaitRec struct {
	ContChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (r WaitRec) Via() symbol.ADT { return r.ContChnlPH }

func (WaitRec) impl() {}

type SendRec struct {
	ConnChnlPH symbol.ADT
	ContChnlID identity.ADT
	ValChnlID  identity.ADT
	ValExpCK   contkey.ADT
	ConnID     identity.ADT
}

func (r SendRec) Via() symbol.ADT { return r.ConnChnlPH }

func (SendRec) impl() {}

type RecvRec struct {
	ConnChnlPH symbol.ADT
	ContChnlID identity.ADT
	NewChnlPH  symbol.ADT
	ContExp    ExpSpec
}

func (r RecvRec) Via() symbol.ADT { return r.ConnChnlPH }

func (RecvRec) impl() {}

type LabRec struct {
	ConnChnlPH symbol.ADT
	ContChnlID identity.ADT
	ValExpCK   contkey.ADT
	ValLabQN   uniqsym.ADT // TODO выпилить
}

func (r LabRec) Via() symbol.ADT { return r.ConnChnlPH }

func (LabRec) impl() {}

type CaseRec struct {
	ConnChnlPH symbol.ADT
	ContChnlID identity.ADT
	ContExps   map[uniqsym.ADT]ExpSpec
}

func (r CaseRec) Via() symbol.ADT { return r.ConnChnlPH }

func (CaseRec) impl() {}

type FwdRec struct {
	ConnChnlPH symbol.ADT
	ContChnlID identity.ADT
}

func (r FwdRec) Via() symbol.ADT { return r.ConnChnlPH }

func (FwdRec) impl() {}

func CollectEnv(spec ExpSpec) []identity.ADT {
	return collectEnvRec(spec, []identity.ADT{})
}

func collectEnvRec(s ExpSpec, env []identity.ADT) []identity.ADT {
	switch spec := s.(type) {
	case RecvSpec:
		return collectEnvRec(spec.ContExp, env)
	case CaseSpec:
		for _, cont := range spec.ContExps {
			env = collectEnvRec(cont, env)
		}
		return env
	default:
		return env
	}
}

func ErrExpTypeUnexpected(got ExpSpec) error {
	return fmt.Errorf("exp spec unexpected: %T", got)
}

func ErrRecTypeUnexpected(got ExpRec) error {
	return fmt.Errorf("exp rec unexpected: %T", got)
}

func ErrExpTypeMismatch(got, want ExpSpec) error {
	return fmt.Errorf("exp spec mismatch: want %T, got %T", want, got)
}

func ErrExpValueNil(pid identity.ADT) error {
	return fmt.Errorf("proc %q term is nil", pid)
}
