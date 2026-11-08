package core

import (
	"fmt"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

type API interface {
	Create(DefSpec) (typesem.SemRef, error)
	Retrieve(identity.ADT) (DefRec, error)
}

type DefSpec struct {
	ProcQN uniqsym.ADT // or dec.ProcID
	ProcES termexp.ExpSpec
}

type DefRec struct {
	Ref typesem.SemRef
}

type DefSnap struct {
	Ref typesem.SemRef
}

type service struct {
	procDefs   Repo
	transactor tx.Transactor
	log        *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func NewService(
	procs Repo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	return &service{procs, transactor, log}
}

func (s *service) Create(spec DefSpec) (typesem.SemRef, error) {
	return typesem.SemRef{}, nil
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
