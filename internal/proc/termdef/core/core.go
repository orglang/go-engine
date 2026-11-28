package core

import (
	"context"
	"fmt"

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
	Expr    ExprSpec
}

type DefRec struct {
	TermID identity.ADT
	ExprCK  contkey.ADT
}

type service struct {
	termDecRepo termdec.Repo
	termDefRepo Repo
	transactor  tx.Transactor
}

func NewService(termDecRepo termdec.Repo, termDefRepo Repo, transactor tx.Transactor) API {
	return &service{
		termDecRepo: termDecRepo,
		termDefRepo: termDefRepo,
		transactor:  transactor,
	}
}

var _ API = (*service)(nil)

func (s *service) RetrieveDefinition(ctx context.Context, termID identity.ADT) (DefSpec, error) {
	return tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DefSpec, error) {
		rec, err := s.termDefRepo.GetRecByID(txCtx, conn, termID)
		if err != nil {
			return DefSpec{}, err
		}
		if rec.ExprCK == 0 {
			return DefSpec{}, DefinitionMissing(termID)
		}
		termQN, err := s.termDecRepo.GetQNByID(txCtx, conn, rec.TermID)
		if err != nil {
			return DefSpec{}, err
		}
		expr, err := s.termDefRepo.SelectSpecByCK(txCtx, conn, rec.ExprCK)
		if err != nil {
			return DefSpec{}, err
		}
		return DefSpec{TermQN: termQN, Expr: expr}, nil
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
