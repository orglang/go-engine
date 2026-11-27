package core

import (
	"context"
	"log/slog"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
	typedef "orglang/go-engine/internal/pool/typedef/core"
)

type service struct {
	repo         Repo
	descRootRepo descroot.Repo
	typeDefRepo  typedef.Repo
	transactor   tx.Transactor
	log          *slog.Logger
}

func NewService(repo Repo, descRootRepo descroot.Repo, typeDefRepo typedef.Repo, transactor tx.Transactor, log *slog.Logger) API {
	return &service{repo: repo, descRootRepo: descRootRepo, typeDefRepo: typeDefRepo, transactor: transactor, log: log}
}

var _ API = (*service)(nil)

func (s *service) Create(ctx context.Context, spec DecSpec) (_ termsem.SemRef, err error) {
	qnAttr := slog.Any("qn", spec.TermQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))

	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars)+1)
	for _, varSpec := range spec.AssetVars {
		assetQNs = append(assetQNs, varSpec.TypeQN)
	}

	typeDefs, getErr := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (map[uniqsym.ADT]typedef.DefRec, error) {
		return s.typeDefRepo.GetRecsByQNs(txCtx, conn, append(assetQNs, spec.LiabVar.TypeQN))
	})
	if getErr != nil {
		return termsem.SemRef{}, getErr
	}

	liabTypeDef := typeDefs[spec.LiabVar.TypeQN]
	newLiabVar := termvar.VarRec{TypeID: liabTypeDef.TypeID, ChnlPH: spec.LiabVar.ChnlPH, ExpCK: liabTypeDef.ExpCK}
	newAssetVars := make([]termvar.VarRec, 0, len(spec.AssetVars))
	for _, varSpec := range spec.AssetVars {
		typeDef := typeDefs[varSpec.TypeQN]
		newAssetVars = append(newAssetVars, termvar.VarRec{TypeID: typeDef.TypeID, ChnlPH: varSpec.ChnlPH, ExpCK: typeDef.ExpCK})
	}

	expectedRoot := liabTypeDef.DescRef
	newRoot := descroot.SemRef{DescID: expectedRoot.DescID, DescRN: seqnum.Next(expectedRoot.DescRN)}
	newDec := DecRec{TermQN: spec.TermQN, TermRef: termsem.New(newRoot), LiabVar: newLiabVar, AssetVars: newAssetVars}

	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if err := s.descRootRepo.TouchRef(txCtx, conn, expectedRoot); err != nil { return err }
		return s.repo.AddRec(txCtx, conn, newDec)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return termsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newDec.TermRef))
	return newDec.TermRef, nil
}
