package core

import (
	"context"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
	termdec "orglang/go-engine/internal/pool/termdec/core"
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
	Exp    ExpSpec
}

func NewService(termDecRepo termdec.Repo, termDefRepo Repo, transactor tx.Transactor) API {
	return &service{termDecRepo: termDecRepo, termDefRepo: termDefRepo, transactor: transactor}
}

type service struct {
	termDecRepo termdec.Repo
	termDefRepo Repo
	transactor  tx.Transactor
}

var _ API = (*service)(nil)

func (s *service) RetrieveDefinition(ctx context.Context, termID identity.ADT) (DefSpec, error) {
	return tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DefSpec, error) {
		rec, err := s.termDefRepo.GetByTermID(txCtx, conn, termID)
		if err != nil {
			return DefSpec{}, err
		}
		if rec.Exp == nil {
			return DefSpec{}, DefinitionMissing(termID)
		}
		termQN, err := s.termDecRepo.GetQNByID(txCtx, conn, rec.TermID)
		if err != nil {
			return DefSpec{}, err
		}
		return DefSpec{TermQN: termQN, Exp: rec.Exp}, nil
	})
}
