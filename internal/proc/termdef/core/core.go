package core

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/termvar"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	"orglang/go-engine/internal/adt/uniqsym"

	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type API interface {
	Incept(context.Context, uniqsym.ADT) (termsem.SemRef, error)
	Create(context.Context, DecSpec) (DecSnap, error)
	RetrieveSnap(context.Context, termsem.SemRef) (DecSnap, error)
	RetrieveDefinition(context.Context, identity.ADT) (DefSpec, error)
	RetreiveRefs(context.Context) ([]termsem.SemRef, error)
}

type DefSpec struct {
	TermQN uniqsym.ADT
	Exp    ExpSpec
}

type DecSpec = termdec.DecSpec

type DecRec = termdec.DecRec

type DefRec struct {
	TermID identity.ADT
	ExpCK  contkey.ADT
}

// aka ExpDec or ExpDecDef without expression
type DecSnap = termdec.DecSnap

type service struct {
	termDecRepo  termdec.Repo
	termDefRepo  Repo
	descRootRepo descroot.Repo
	descSemRepo  descsem.Repo
	typeDefRepo  typedef.Repo
	transactor   tx.Transactor
	log          *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func NewService(
	termDecRepo termdec.Repo,
	termDefRepo Repo,
	descRootRepo descroot.Repo,
	descSemRepo descsem.Repo,
	typeDefRepo typedef.Repo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	return &service{termDecRepo: termDecRepo, termDefRepo: termDefRepo, descRootRepo: descRootRepo, descSemRepo: descSemRepo, typeDefRepo: typeDefRepo, transactor: transactor, log: log}
}

func (s *service) Incept(ctx context.Context, termQN uniqsym.ADT) (_ termsem.SemRef, err error) {
	qnAttr := slog.Any("qn", termQN)
	s.log.Debug("inception started", qnAttr)
	newRoot := descroot.New()
	newDec := DecRec{TermRef: termsem.New(newRoot), TermQN: termQN}
	newBind := descsem.SemRec{DescQN: termQN, DescID: newDec.TermRef.TermID, Kind: descsem.TermKind}
	err = s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if err := s.descRootRepo.AddRef(txCtx, conn, newRoot); err != nil {
			return err
		}
		if err := s.descSemRepo.AddRec(txCtx, conn, newBind); err != nil {
			return err
		}
		return s.termDecRepo.AddRec(txCtx, conn, newDec)
	})
	if err != nil {
		s.log.Error("inception failed", qnAttr)
		return termsem.SemRef{}, err
	}
	s.log.Debug("inception succeed", qnAttr, slog.Any("ref", newDec.TermRef))
	return newDec.TermRef, nil
}

func (s *service) Create(ctx context.Context, spec DecSpec) (_ DecSnap, err error) {
	qnAttr := slog.Any("qn", spec.TermQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars)+1)
	for _, spec := range spec.AssetVars {
		assetQNs = append(assetQNs, spec.TypeQN)
	}
	typeDefs, getErr := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (map[uniqsym.ADT]typedef.DefRec, error) {
		return s.typeDefRepo.SelectEnv(txCtx, conn, append(assetQNs, spec.LiabVar.TypeQN))
	})
	if getErr != nil {
		return DecSnap{}, getErr
	}
	providerDef := typeDefs[spec.LiabVar.TypeQN]
	newLiabVar := termvar.VarRec{
		TypeRef: providerDef.TypeRef,
		ChnlPH:  spec.LiabVar.ChnlPH,
		ExpCK:   providerDef.ExpCK,
	}
	newAssetVars := make([]termvar.VarRec, 0, len(spec.AssetVars))
	for _, assetVar := range spec.AssetVars {
		typeDef := typeDefs[assetVar.TypeQN]
		newAssetVars = append(newAssetVars, termvar.VarRec{
			TypeRef: typeDef.TypeRef,
			ChnlPH:  assetVar.ChnlPH,
			ExpCK:   typeDef.ExpCK,
		})
	}
	providerRef := providerDef.TypeRef
	expectedRoot := descroot.SemRef{DescID: providerRef.DescID, DescRN: providerRef.DescRN}
	newRoot := descroot.SemRef{DescID: expectedRoot.DescID, DescRN: seqnum.Next(expectedRoot.DescRN)}
	newDec := DecRec{TermRef: termsem.New(newRoot), TermQN: spec.TermQN, LiabVar: newLiabVar, AssetVars: newAssetVars}
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
		return DecSnap{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newDec.TermRef))
	return DecSnap{TermRef: newDec.TermRef, LiabVar: newLiabVar, AssetVars: newAssetVars}, nil
}

func (s *service) RetrieveSnap(ctx context.Context, ref termsem.SemRef) (snap DecSnap, err error) {
	snap, err = tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DecSnap, error) {
		return s.termDecRepo.GetSnap(txCtx, conn, ref)
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("ref", ref))
		return DecSnap{}, err
	}
	return snap, nil
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
		termQN, err := s.termDefRepo.GetQNByID(txCtx, conn, rec.TermID)
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

func (s *service) RetreiveRefs(ctx context.Context) (refs []termsem.SemRef, err error) {
	refs, err = tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) ([]termsem.SemRef, error) {
		return s.termDecRepo.GetRefs(txCtx, conn)
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
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
