package procdef

import (
	"fmt"
	"log/slog"

	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
)

type API interface {
	Create(DefSpec) (DefRef, error)
	Retrieve(identity.ADT) (DefRec, error)
}

type DefSpec struct {
	ProcQN qualsym.ADT // or dec.ProcID
	ProcTS TermSpec
}

type DefRef struct {
	DefID identity.ADT
}

type DefRec struct {
	DefID identity.ADT
}

type DefSnap struct {
	DefID identity.ADT
}

type TermSpec interface {
	Via() qualsym.ADT
}

type CloseSpec struct {
	CommPH qualsym.ADT
}

func (s CloseSpec) Via() qualsym.ADT { return s.CommPH }

type WaitSpec struct {
	CommPH qualsym.ADT
	ContES TermSpec
}

func (s WaitSpec) Via() qualsym.ADT { return s.CommPH }

type SendSpec struct {
	CommPH qualsym.ADT // via
	ValPH  qualsym.ADT // val
}

func (s SendSpec) Via() qualsym.ADT { return s.CommPH }

type RecvSpec struct {
	CommPH qualsym.ADT
	BindPH qualsym.ADT
	ContES TermSpec
}

func (s RecvSpec) Via() qualsym.ADT { return s.CommPH }

type LabSpec struct {
	CommPH qualsym.ADT
	Label  qualsym.ADT
	ContES TermSpec
}

func (s LabSpec) Via() qualsym.ADT { return s.CommPH }

type CaseSpec struct {
	CommPH  qualsym.ADT
	ContESs map[qualsym.ADT]TermSpec
}

func (s CaseSpec) Via() qualsym.ADT { return s.CommPH }

// aka ExpName
type LinkSpec struct {
	ProcQN qualsym.ADT
	X      identity.ADT
	Ys     []identity.ADT
}

func (s LinkSpec) Via() qualsym.ADT { return "" }

type FwdSpec struct {
	X qualsym.ADT // old via (from)
	Y qualsym.ADT // new via (to)
}

func (s FwdSpec) Via() qualsym.ADT { return s.X }

// аналог SendSpec, но значения отправляются балком
type CallSpecOld struct {
	X     qualsym.ADT
	SigPH qualsym.ADT // import
	Ys    []qualsym.ADT
}

func (s CallSpecOld) Via() qualsym.ADT { return s.SigPH }

type CallSpec struct {
	CommPH qualsym.ADT
	BindPH qualsym.ADT
	ProcSN qualsym.ADT
	ValPHs []qualsym.ADT // channel bulk
	ContES TermSpec
}

func (s CallSpec) Via() qualsym.ADT { return s.CommPH }

// аналог RecvSpec, но значения принимаются балком
type SpawnSpecOld struct {
	X      qualsym.ADT
	SigID  identity.ADT
	Ys     []qualsym.ADT
	PoolQN qualsym.ADT
	ContES TermSpec
}

func (s SpawnSpecOld) Via() qualsym.ADT { return s.X }

type SpawnSpec struct {
	CommPH qualsym.ADT
	ProcSN qualsym.ADT
	ContES TermSpec
}

func (s SpawnSpec) Via() qualsym.ADT { return s.CommPH }

type AcqureSpec struct {
	CommPH qualsym.ADT
	ContES TermSpec
}

func (s AcqureSpec) Via() qualsym.ADT { return s.CommPH }

type AcceptSpec struct {
	CommPH qualsym.ADT
	ContES TermSpec
}

func (s AcceptSpec) Via() qualsym.ADT { return s.CommPH }

type DetachSpec struct {
	CommPH qualsym.ADT
}

func (s DetachSpec) Via() qualsym.ADT { return s.CommPH }

type ReleaseSpec struct {
	CommPH qualsym.ADT
}

func (s ReleaseSpec) Via() qualsym.ADT { return s.CommPH }

type TermRec interface {
	TermSpec
	impl()
}

type CloseRec struct {
	X qualsym.ADT
}

func (r CloseRec) Via() qualsym.ADT { return r.X }

func (CloseRec) impl() {}

type WaitRec struct {
	X      qualsym.ADT
	ContES TermSpec
}

func (r WaitRec) Via() qualsym.ADT { return r.X }

func (WaitRec) impl() {}

type SendRec struct {
	X     qualsym.ADT
	A     identity.ADT
	B     identity.ADT
	ExpID identity.ADT
}

func (r SendRec) Via() qualsym.ADT { return r.X }

func (SendRec) impl() {}

type RecvRec struct {
	X      qualsym.ADT
	A      identity.ADT
	Y      qualsym.ADT
	ContES TermSpec
}

func (r RecvRec) Via() qualsym.ADT { return r.X }

func (RecvRec) impl() {}

type LabRec struct {
	X     qualsym.ADT
	A     identity.ADT
	Label qualsym.ADT
}

func (r LabRec) Via() qualsym.ADT { return r.X }

func (LabRec) impl() {}

type CaseRec struct {
	X       qualsym.ADT
	A       identity.ADT
	ContESs map[qualsym.ADT]TermSpec
}

func (r CaseRec) Via() qualsym.ADT { return r.X }

func (CaseRec) impl() {}

type FwdRec struct {
	X qualsym.ADT
	B identity.ADT // to
}

func (r FwdRec) Via() qualsym.ADT { return r.X }

func (FwdRec) impl() {}

func CollectEnv(spec TermSpec) []identity.ADT {
	return collectEnvRec(spec, []identity.ADT{})
}

type service struct {
	procs    Repo
	operator sd.Operator
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

func newService(
	procs Repo,
	operator sd.Operator,
	l *slog.Logger,
) *service {
	return &service{procs, operator, l}
}

func (s *service) Create(spec DefSpec) (DefRef, error) {
	return DefRef{}, nil
}

func (s *service) Retrieve(recID identity.ADT) (DefRec, error) {
	return DefRec{}, nil
}

func collectEnvRec(s TermSpec, env []identity.ADT) []identity.ADT {
	switch spec := s.(type) {
	case RecvSpec:
		return collectEnvRec(spec.ContES, env)
	case CaseSpec:
		for _, cont := range spec.ContESs {
			env = collectEnvRec(cont, env)
		}
		return env
	case SpawnSpecOld:
		return collectEnvRec(spec.ContES, append(env, spec.SigID))
	default:
		return env
	}
}

func ErrDoesNotExist(want identity.ADT) error {
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

func ErrTermValueNil(pid identity.ADT) error {
	return fmt.Errorf("proc %q term is nil", pid)
}

func ErrMissingInCfg(want qualsym.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func ErrMissingInCfg2(want identity.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func ErrMissingInCtx(want qualsym.ADT) error {
	return fmt.Errorf("channel missing in ctx: %v", want)
}
