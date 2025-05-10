package def

import (
	"fmt"

	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type TermSpec interface {
	Via() sym.ADT
}

type CloseSpec struct {
	X sym.ADT
}

func (s CloseSpec) Via() sym.ADT { return s.X }

type WaitSpec struct {
	X    sym.ADT
	Cont TermSpec
}

func (s WaitSpec) Via() sym.ADT { return s.X }

type SendSpec struct {
	X sym.ADT // via
	Y sym.ADT // val
}

func (s SendSpec) Via() sym.ADT { return s.X }

type RecvSpec struct {
	X    sym.ADT // via
	Y    sym.ADT // val
	Cont TermSpec
}

func (s RecvSpec) Via() sym.ADT { return s.X }

type LabSpec struct {
	X     sym.ADT
	Label sym.ADT
}

func (s LabSpec) Via() sym.ADT { return s.X }

type CaseSpec struct {
	X     sym.ADT
	Conts map[sym.ADT]TermSpec
}

func (s CaseSpec) Via() sym.ADT { return s.X }

// aka ExpName
type LinkSpec struct {
	SigQN sym.ADT
	X     id.ADT
	Ys    []id.ADT
}

func (s LinkSpec) Via() sym.ADT { return "" }

type FwdSpec struct {
	X sym.ADT // old via (from)
	Y sym.ADT // new via (to)
}

func (s FwdSpec) Via() sym.ADT { return s.X }

// аналог SendSpec, но значения отправляются балком
type CallSpec struct {
	X     sym.ADT
	SigPH sym.ADT // import
	Ys    []sym.ADT
}

func (s CallSpec) Via() sym.ADT { return s.SigPH }

// аналог RecvSpec, но значения принимаются балком
type SpawnSpec struct {
	X      sym.ADT
	SigID  id.ADT
	Ys     []sym.ADT
	PoolQN sym.ADT
	Cont   TermSpec
}

func (s SpawnSpec) Via() sym.ADT { return s.X }

// аналог RecvSpec, но значения принимаются балком
type SpawnSpec2 struct {
	X     sym.ADT
	SigPH sym.ADT // export
	Cont  TermSpec
}

func (s SpawnSpec2) Via() sym.ADT { return s.SigPH }

type TermRec interface {
	TermSpec
	impl()
}

type CloseRec struct {
	X sym.ADT
}

func (r CloseRec) Via() sym.ADT { return r.X }

func (CloseRec) impl() {}

type WaitRec struct {
	X    sym.ADT
	Cont TermSpec
}

func (r WaitRec) Via() sym.ADT { return r.X }

func (WaitRec) impl() {}

type SendRec struct {
	X      sym.ADT
	A      id.ADT
	B      id.ADT
	TermID id.ADT
}

func (r SendRec) Via() sym.ADT { return r.X }

func (SendRec) impl() {}

type RecvRec struct {
	X    sym.ADT
	A    id.ADT
	Y    sym.ADT
	Cont TermSpec
}

func (r RecvRec) Via() sym.ADT { return r.X }

func (RecvRec) impl() {}

type LabRec struct {
	X     sym.ADT
	A     id.ADT
	Label sym.ADT
}

func (r LabRec) Via() sym.ADT { return r.X }

func (LabRec) impl() {}

type CaseRec struct {
	X     sym.ADT
	A     id.ADT
	Conts map[sym.ADT]TermSpec
}

func (r CaseRec) Via() sym.ADT { return r.X }

func (CaseRec) impl() {}

type FwdRec struct {
	X sym.ADT
	B id.ADT // to
}

func (r FwdRec) Via() sym.ADT { return r.X }

func (FwdRec) impl() {}

func CollectEnv(spec TermSpec) []id.ADT {
	return collectEnvRec(spec, []id.ADT{})
}

type Repo interface {
}

func collectEnvRec(s TermSpec, env []id.ADT) []id.ADT {
	switch spec := s.(type) {
	case RecvSpec:
		return collectEnvRec(spec.Cont, env)
	case CaseSpec:
		for _, cont := range spec.Conts {
			env = collectEnvRec(cont, env)
		}
		return env
	case SpawnSpec:
		return collectEnvRec(spec.Cont, append(env, spec.SigID))
	default:
		return env
	}
}

func ErrDoesNotExist(want id.ADT) error {
	return fmt.Errorf("rec doesn't exist: %v", want)
}

func ErrTermTypeUnexpected(got TermSpec) error {
	return fmt.Errorf("term spec unexpected: %T", got)
}

func ErrRecTypeUnexpected(got TermRec) error {
	return fmt.Errorf("term rec unexpected: %T", got)
}

func ErrTermTypeMismatch(got, want TermSpec) error {
	return fmt.Errorf("term spec mismatch: want %T, got %T", want, got)
}

func ErrTermValueNil(pid id.ADT) error {
	return fmt.Errorf("proc %q term is nil", pid)
}

func ErrMissingInCfg(want sym.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func ErrMissingInCfg2(want id.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func ErrMissingInCtx(want sym.ADT) error {
	return fmt.Errorf("channel missing in ctx: %v", want)
}
