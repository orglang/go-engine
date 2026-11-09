package core

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
	typeexp "orglang/go-engine/internal/proc/typeexp/core"
)

type API interface {
	Create(DefSpec) (DefSnap, error)
	Modify(DefSnap) (DefSnap, error)
	RetrieveSnap(typesem.SemRef) (DefSnap, error)
	retrieveSnap(DefRec) (DefSnap, error)
	RetreiveRefs() ([]typesem.SemRef, error)
}

type DefSpec struct {
	TypeQN  uniqsym.ADT
	TypeExp typeexp.ExpSpec
}

// aka TpDef
type DefRec struct {
	TypeRef typesem.SemRef
	ExpCK   contkey.ADT
}

type DefSnap struct {
	TypeRef typesem.SemRef
	DefSpec DefSpec
}

type Context struct {
	Assets map[symbol.ADT]typeexp.ExpRec
	Liabs  map[symbol.ADT]typeexp.ExpRec
}

type service struct {
	typeDefRepo Repo
	typeExpRepo typeexp.Repo
	descSemRepo descsem.Repo
	transactor  tx.Transactor
	log         *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func NewService(
	typeDefRepo Repo,
	typeExpRepo typeexp.Repo,
	descSemRepo descsem.Repo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	return &service{typeDefRepo, typeExpRepo, descSemRepo, transactor, log}
}

func (s *service) Create(spec DefSpec) (_ DefSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", spec.TypeQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newExp, err := typeexp.ConvertSpecToRec(spec.TypeExp)
	if err != nil {
		return DefSnap{}, err
	}
	newDef := DefRec{TypeRef: typesem.New(), ExpCK: newExp.Key()}
	newDesc := descsem.SemRec{DescQN: spec.TypeQN, DescID: newDef.TypeRef.TypeID, Kind: descsem.TypeKind}
	err = s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		err = s.descSemRepo.AddRec(txCtx, conn, newDesc)
		if err != nil {
			return err
		}
		err = s.typeExpRepo.AddRec(txCtx, conn, newExp)
		if err != nil {
			return err
		}
		return s.typeDefRepo.AddRec(txCtx, conn, newDef)
	})
	if err != nil {
		s.log.Error("creation failed", qnAttr)
		return DefSnap{}, err
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newDef.TypeRef))
	return DefSnap{
		TypeRef: newDef.TypeRef,
		DefSpec: spec,
	}, nil
}

func (s *service) Modify(snap DefSnap) (_ DefSnap, err error) {
	ctx := context.Background()
	refAttr := slog.Any("defRef", snap.TypeRef)
	s.log.Debug("modification started", refAttr)
	rec, err := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DefRec, error) {
		return s.typeDefRepo.GetRecByRef(txCtx, conn, snap.TypeRef)
	})
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	if snap.TypeRef.TypeRN != rec.TypeRef.TypeRN {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, errConcurrentModification(snap.TypeRef.TypeRN, rec.TypeRef.TypeRN)
	}
	snap.TypeRef.TypeRN = seqnum.Next(snap.TypeRef.TypeRN)
	curSnap, err := s.retrieveSnap(rec)
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	err = s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if typeexp.CheckSpec(snap.DefSpec.TypeExp, curSnap.DefSpec.TypeExp) != nil {
			newExp, err := typeexp.ConvertSpecToRec(snap.DefSpec.TypeExp)
			if err != nil {
				return err
			}
			err = s.typeExpRepo.AddRec(txCtx, conn, newExp)
			if err != nil {
				return err
			}
			rec.ExpCK = newExp.Key()
			rec.TypeRef.TypeRN = snap.TypeRef.TypeRN
		}
		if rec.TypeRef.TypeRN == snap.TypeRef.TypeRN {
			err = s.typeDefRepo.ModifyRec(txCtx, conn, rec)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	s.log.Debug("modification succeed", refAttr)
	return snap, nil
}

func (s *service) RetrieveSnap(ref typesem.SemRef) (_ DefSnap, err error) {
	ctx := context.Background()
	rec, getErr := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DefRec, error) {
		return s.typeDefRepo.GetRecByRef(txCtx, conn, ref)
	})
	if getErr != nil {
		s.log.Error("retrieval failed", slog.Any("ref", ref))
		return DefSnap{}, getErr
	}
	return s.retrieveSnap(rec)
}

func (s *service) retrieveSnap(rec DefRec) (_ DefSnap, err error) { // revive:disable-line:confusing-naming
	ctx := context.Background()
	expRec, err := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (typeexp.ExpRec, error) {
		return s.typeExpRepo.SelectRecByCK(txCtx, conn, rec.ExpCK)
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("ref", rec.TypeRef))
		return DefSnap{}, err
	}
	return DefSnap{
		TypeRef: rec.TypeRef,
		DefSpec: DefSpec{TypeExp: typeexp.ConvertRecToSpec(expRec)},
	}, nil
}

func (s *service) RetreiveRefs() (refs []typesem.SemRef, err error) {
	ctx := context.Background()
	refs, err = tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) ([]typesem.SemRef, error) {
		return s.typeDefRepo.GetRefs(txCtx, conn)
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}

func CollectEnv(recs iter.Seq[DefRec]) []contkey.ADT {
	expIDs := []contkey.ADT{}
	for r := range recs {
		expIDs = append(expIDs, r.ExpCK)
	}
	return expIDs
}

func ErrSymMissingInEnv(want uniqsym.ADT) error {
	return fmt.Errorf("root missing in env: %v", want)
}

func errConcurrentModification(got seqnum.ADT, want seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: want revision %v, got revision %v", want, got)
}

func ErrOptimisticUpdate(got seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}

func ErrDoesNotExist(want identity.ADT) error {
	return fmt.Errorf("root doesn't exist: %v", want)
}

func ErrMissingInEnv(want contkey.ADT) error {
	return fmt.Errorf("exp missing in env: %v", want)
}

func ErrMissingInCfg(want identity.ADT) error {
	return fmt.Errorf("root missing in cfg: %v", want)
}

func ErrMissingInCtx(want symbol.ADT) error {
	return fmt.Errorf("root missing in ctx: %v", want)
}
