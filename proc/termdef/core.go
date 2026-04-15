package termdef

import (
	"fmt"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/semtype"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/proc/termexp"
)

type API interface {
	Create(DefSpec) (semtype.TypeRef, error)
	Retrieve(identity.ADT) (DefRec, error)
}

type DefSpec struct {
	ProcQN uniqsym.ADT // or dec.ProcID
	ProcES termexp.ExpSpec
}

type DefRec struct {
	Ref semtype.TypeRef
}

type DefSnap struct {
	Ref semtype.TypeRef
}

type service struct {
	procDefs Repo
	operator db.Operator
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	procs Repo,
	operator db.Operator,
	l *slog.Logger,
) *service {
	return &service{procs, operator, l}
}

func (s *service) Create(spec DefSpec) (semtype.TypeRef, error) {
	return semtype.TypeRef{}, nil
}

func (s *service) Retrieve(recID identity.ADT) (DefRec, error) {
	return DefRec{}, nil
}

func ErrDoesNotExist(want identity.ADT) error {
	return fmt.Errorf("rec doesn't exist: %v", want)
}

func ErrMissingInCfg(want symbol.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func ErrMissingInCfg2(want identity.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func ErrMissingInCtx(want symbol.ADT) error {
	return fmt.Errorf("channel missing in ctx: %v", want)
}
