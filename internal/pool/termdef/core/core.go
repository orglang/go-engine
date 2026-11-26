package core

import (
	"context"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/uniqsym"
	typedef "orglang/go-engine/internal/pool/typedef/core"
	termdec "orglang/go-engine/internal/pool/termdec/core"
)

type API interface {
	Create(context.Context, termdec.DecSpec) (termsem.SemRef, error)
	RetrieveDefinition(context.Context, identity.ADT) (DefSpec, error)
}

type DefSpec struct {
	TermQN uniqsym.ADT
	Exp    ExpSpec
}

type DecRec = termdec.DecRec

type DefRec struct {
	TermID identity.ADT
	ExpCK  contkey.ADT
}

func NewService(
	termDecRepo termdec.Repo,
	termDefRepo Repo,
	descRootRepo descroot.Repo,
	typeDefRepo typedef.Repo,
	descSemRepo descsem.Repo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	return &service{termDecRepo: termDecRepo, termDefRepo: termDefRepo, descRootRepo: descRootRepo, typeDefRepo: typeDefRepo, descSemRepo: descSemRepo, transactor: transactor, log: log}
}

type service struct {
	termDecRepo  termdec.Repo
	termDefRepo  Repo
	descRootRepo descroot.Repo
	typeDefRepo  typedef.Repo
	descSemRepo  descsem.Repo
	transactor   tx.Transactor
	log          *slog.Logger
}

func (s *service) Create(ctx context.Context, spec termdec.DecSpec) (_ termsem.SemRef, err error) {
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
	newLiabVar := termvar.VarRec{
		TypeRef: typeDefs[spec.LiabVar.TypeQN].TypeRef,
		ChnlPH:  spec.LiabVar.ChnlPH,
		ExpCK:   typeDefs[spec.LiabVar.TypeQN].ExpCK,
	}
	newAssetVars := make([]termvar.VarRec, 0, len(spec.AssetVars))
	for _, varSpec := range spec.AssetVars {
		newAssetVars = append(newAssetVars, termvar.VarRec{
			TypeRef: typeDefs[varSpec.TypeQN].TypeRef,
			ChnlPH:  varSpec.ChnlPH,
			ExpCK:   typeDefs[varSpec.TypeQN].ExpCK,
		})
	}
	providerRef := typeDefs[spec.LiabVar.TypeQN].TypeRef
	expectedRoot := descroot.SemRef{DescID: providerRef.DescID, DescRN: providerRef.DescRN}
	newRoot := descroot.SemRef{DescID: expectedRoot.DescID, DescRN: seqnum.Next(expectedRoot.DescRN)}
	newDec := DecRec{TermRef: termsem.New(newRoot), LiabVar: newLiabVar, AssetVars: newAssetVars}
	newBind := descsem.SemRec{DescQN: spec.TermQN, DescID: newDec.TermRef.TermID, Kind: descsem.TermKind}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if err := s.descRootRepo.TouchRef(txCtx, conn, expectedRoot); err != nil {
			return err
		}
		if err := s.descSemRepo.AddRec(txCtx, conn, newBind); err != nil {
			return err
		}
		return s.termDecRepo.AddRec(txCtx, conn, newDec)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return termsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newDec.TermRef))
	return newDec.TermRef, nil
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
