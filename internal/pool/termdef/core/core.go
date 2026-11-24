package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/uniqsym"
	termexp "orglang/go-engine/internal/pool/termexp/core"
	typedef "orglang/go-engine/internal/pool/typedef/core"
)

type API interface {
	Create(context.Context, DefSpec) (termsem.SemRef, error)
	RetrieveDefinition(context.Context, identity.ADT, *seqnum.ADT) (termexp.DescSpec, error)
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

type DefSpec struct {
	TermQN    uniqsym.ADT
	LiabVar   termvar.VarSpec
	AssetVars []termvar.VarSpec
}

type DefRec struct {
	TermRef   termsem.SemRef
	ExpCK     contkey.ADT
	LiabVar   termvar.VarRec
	AssetVars []termvar.VarRec
}

type DefSnap struct {
	TermRef termsem.SemRef
}

func NewService(
	termDefRepo Repo,
	termExpRepo termexp.Repo,
	descRootRepo descroot.Repo,
	typeDefRepo typedef.Repo,
	descSemRepo descsem.Repo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	return &service{termDefRepo: termDefRepo, termExpRepo: termExpRepo, descRootRepo: descRootRepo, typeDefRepo: typeDefRepo, descSemRepo: descSemRepo, transactor: transactor, log: log}
}

type service struct {
	termDefRepo  Repo
	termExpRepo  termexp.Repo
	descRootRepo descroot.Repo
	typeDefRepo  typedef.Repo
	descSemRepo  descsem.Repo
	transactor   tx.Transactor
	log          *slog.Logger
}

func (s *service) Create(ctx context.Context, spec DefSpec) (_ termsem.SemRef, err error) {
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
	newDef := DefRec{TermRef: termsem.New(newRoot), LiabVar: newLiabVar, AssetVars: newAssetVars}
	newBind := descsem.SemRec{DescQN: spec.TermQN, DescID: newDef.TermRef.TermID, Kind: descsem.TermKind}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		err = s.descRootRepo.TouchRef(txCtx, conn, expectedRoot)
		if err != nil {
			return err
		}
		err = s.descSemRepo.AddRec(txCtx, conn, newBind)
		if err != nil {
			return err
		}
		return s.termDefRepo.AddRec(txCtx, conn, newDef)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return termsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newDef.TermRef))
	return newDef.TermRef, nil
}

func (s *service) RetrieveDefinition(ctx context.Context, termID identity.ADT, expectedRN *seqnum.ADT) (termexp.DescSpec, error) {
	if s.termExpRepo == nil {
		return nil, fmt.Errorf("term Definition repository is not configured")
	}
	return tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (termexp.DescSpec, error) {
		rec, err := s.termDefRepo.GetRecByID(txCtx, conn, termID)
		if err != nil {
			return nil, err
		}
		if expectedRN != nil && *expectedRN != rec.TermRef.DescRN {
			return nil, fmt.Errorf("%w: expected %d, current %d", ErrDefinitionRevisionMismatch, *expectedRN, rec.TermRef.DescRN)
		}
		if rec.ExpCK == 0 {
			return nil, ErrDefinitionMissing(termID)
		}
		return s.termExpRepo.SelectSpecByCK(txCtx, conn, rec.ExpCK)
	})
}

var ErrDefinitionRevisionMismatch = errors.New("description revision mismatch")

type DefinitionMissingError struct {
	TermID identity.ADT
}

func (e DefinitionMissingError) Error() string {
	return fmt.Sprintf("pool term %v has no Definition", e.TermID)
}

func ErrDefinitionMissing(termID identity.ADT) error {
	return DefinitionMissingError{TermID: termID}
}
