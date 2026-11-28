package core

import (
	"context"
	"log/slog"

	"orglang/go-engine/internal/adt/progroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type service struct {
	termDecRepo  Repo
	progRootRepo progroot.Repo
	typeDefRepo  typedef.Repo
	transactor   tx.Transactor
	log          *slog.Logger
}

func NewService(termDecRepo Repo, progRootRepo progroot.Repo, typeDefRepo typedef.Repo, transactor tx.Transactor, log *slog.Logger) API {
	return &service{termDecRepo: termDecRepo, progRootRepo: progRootRepo, typeDefRepo: typeDefRepo, transactor: transactor, log: log}
}

func (s *service) Incept(ctx context.Context, termQN uniqsym.ADT) (_ identity.ADT, _ progroot.RootRef, err error) {
	qnAttr := slog.Any("qn", termQN)
	s.log.Debug("inception started", qnAttr)
	termID, newRoot := identity.New(), progroot.New()
	newDec := DecRec{TermID: termID, ProgRef: newRoot, TermQN: termQN}
	err = s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if err := s.progRootRepo.AddRef(txCtx, conn, newRoot); err != nil {
			return err
		}
		return s.termDecRepo.AddRec(txCtx, conn, newDec)
	})
	if err != nil {
		s.log.Error("inception failed", qnAttr)
		return identity.ADT{}, progroot.RootRef{}, err
	}
	s.log.Debug("inception succeed", qnAttr, slog.Any("termID", termID), slog.Any("progRef", newRoot))
	return termID, newRoot, nil
}

func (s *service) Create(ctx context.Context, spec DecSpec) (_ DecSnap, err error) {
	qnAttr := slog.Any("qn", spec.TermQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars)+1)
	for _, assetVar := range spec.AssetVars {
		assetQNs = append(assetQNs, assetVar.TypeQN)
	}
	typeDefs, getErr := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (map[uniqsym.ADT]typedef.DefRec, error) {
		return s.typeDefRepo.SelectEnv(txCtx, conn, append(assetQNs, spec.LiabVar.TypeQN))
	})
	if getErr != nil {
		return DecSnap{}, getErr
	}
	providerDef := typeDefs[spec.LiabVar.TypeQN]
	newLiabVar := termvar.VarRec{TypeID: providerDef.TypeID, ChnlPH: spec.LiabVar.ChnlPH, ExpCK: providerDef.ExpCK}
	newAssetVars := make([]termvar.VarRec, 0, len(spec.AssetVars))
	for _, assetVar := range spec.AssetVars {
		typeDef := typeDefs[assetVar.TypeQN]
		newAssetVars = append(newAssetVars, termvar.VarRec{TypeID: typeDef.TypeID, ChnlPH: assetVar.ChnlPH, ExpCK: typeDef.ExpCK})
	}
	expectedRoot := providerDef.ProgRef
	newRoot := progroot.RootRef{ProgID: expectedRoot.ProgID, ProgRN: seqnum.Next(expectedRoot.ProgRN)}
	termID := identity.New()
	newDec := DecRec{TermID: termID, ProgRef: newRoot, TermQN: spec.TermQN, LiabVar: newLiabVar, AssetVars: newAssetVars}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if err := s.progRootRepo.TouchRef(txCtx, conn, expectedRoot); err != nil {
			return err
		}
		return s.termDecRepo.AddRec(txCtx, conn, newDec)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return DecSnap{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("termID", termID), slog.Any("progRef", newRoot))
	return DecSnap{TermID: termID, ProgRef: newRoot, LiabVar: newLiabVar, AssetVars: newAssetVars}, nil
}

func (s *service) RetrieveSnap(ctx context.Context, id identity.ADT) (snap DecSnap, err error) {
	snap, err = tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DecSnap, error) {
		return s.termDecRepo.GetSnap(txCtx, conn, id)
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("termID", id))
		return DecSnap{}, err
	}
	return snap, nil
}

func (s *service) RetreiveRefs(ctx context.Context) (refs []identity.ADT, err error) {
	refs, err = tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) ([]identity.ADT, error) {
		return s.termDecRepo.GetIDs(txCtx, conn)
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}
