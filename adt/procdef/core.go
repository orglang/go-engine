package procdef

import (
	"fmt"
	"log/slog"

	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/procexp"
	"orglang/orglang/adt/qualsym"
)

type API interface {
	Create(DefSpec) (DefRef, error)
	Retrieve(identity.ADT) (DefRec, error)
}

type DefSpec struct {
	ProcQN qualsym.ADT // or dec.ProcID
	ProcES procexp.ExpSpec
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

func collectEnvRec(s procexp.ExpSpec, env []identity.ADT) []identity.ADT {
	switch spec := s.(type) {
	case procexp.RecvSpec:
		return collectEnvRec(spec.ContES, env)
	case procexp.CaseSpec:
		for _, cont := range spec.ContESs {
			env = collectEnvRec(cont, env)
		}
		return env
	case procexp.SpawnSpecOld:
		return collectEnvRec(spec.ContES, append(env, spec.SigID))
	default:
		return env
	}
}

func ErrDoesNotExist(want identity.ADT) error {
	return fmt.Errorf("rec doesn't exist: %v", want)
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
