package core

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

type API interface {
	RetrieveDefinition(context.Context, identity.ADT) (DefSpec, error)
}

type DefSpec struct {
	TermQN uniqsym.ADT
	Exp    ExpSpec
}

type DefRec struct {
	TermID identity.ADT
	ExpCK  contkey.ADT
}

type service struct {
	termDecRepo termdec.Repo
	termDefRepo Repo
	transactor tx.Transactor
	log *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func NewService(termDecRepo termdec.Repo, termDefRepo Repo, transactor tx.Transactor, log *slog.Logger) API {
	return &service{termDecRepo: termDecRepo, termDefRepo: termDefRepo, transactor: transactor, log: log}
}

func (s *service) RetrieveDefinition(ctx context.Context, termID identity.ADT) (DefSpec, error) {
	return tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DefSpec, error) {
		rec, err := s.termDefRepo.GetRecByID(txCtx, conn, termID)
		if err != nil {
			return DefSpec{}, err
		}
		if rec.ExpCK == 0 {
			return DefSpec{}, DefinitionMissing(termID)
		}
		termQN, err := s.termDecRepo.GetQNByID(txCtx, conn, rec.TermID)
		if err != nil {
			return DefSpec{}, err
		}
		exp, err := s.termDefRepo.SelectSpecByCK(txCtx, conn, rec.ExpCK)
		if err != nil {
			return DefSpec{}, err
		}
		return DefSpec{TermQN: termQN, Exp: exp}, nil
	})
}

func ErrRootMissingInEnv(rid identity.ADT) error {
	return fmt.Errorf("root missing in env: %v", rid)
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
