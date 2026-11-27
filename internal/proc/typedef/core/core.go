package core

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
)

type API interface {
	Create(context.Context, DefSpec) (DefSnap, error)
	Modify(context.Context, DefSnap) (DefSnap, error)
	RetrieveSnap(context.Context, typesem.SemRef) (DefSnap, error)
	retrieveSnap(context.Context, DefRec) (DefSnap, error)
	RetreiveRefs(context.Context) ([]typesem.SemRef, error)
}

type DefSpec struct {
	TypeQN  uniqsym.ADT
	TypeExp ExpSpec
}

// aka TpDef
type DefRec struct {
	TypeQN  uniqsym.ADT
	TypeRef typesem.SemRef
	ExpCK   contkey.ADT
}

type DefSnap struct {
	TypeRef typesem.SemRef
	DefSpec DefSpec
}

type Context struct {
	Assets map[symbol.ADT]ExpRec
	Liabs  map[symbol.ADT]ExpRec
}

type service struct {
	typeDefRepo  Repo
	descRootRepo descroot.Repo
	typeExpRepo  TypeExpRepo
	transactor   tx.Transactor
	log          *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func NewService(
	typeDefRepo Repo,
	descRootRepo descroot.Repo,
	typeExpRepo TypeExpRepo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	return &service{typeDefRepo: typeDefRepo, descRootRepo: descRootRepo, typeExpRepo: typeExpRepo, transactor: transactor, log: log}
}

func (s *service) Create(ctx context.Context, spec DefSpec) (_ DefSnap, err error) {
	qnAttr := slog.Any("qn", spec.TypeQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newExp, err := ConvertSpecToRec(spec.TypeExp)
	if err != nil {
		return DefSnap{}, err
	}
	newDef := DefRec{TypeRef: typesem.New(), TypeQN: spec.TypeQN, ExpCK: newExp.Key()}
	newRoot := descroot.SemRef{DescID: newDef.TypeRef.DescID, DescRN: newDef.TypeRef.DescRN}
	err = s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		err = s.descRootRepo.AddRef(txCtx, conn, newRoot)
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

func (s *service) Modify(ctx context.Context, snap DefSnap) (_ DefSnap, err error) {
	refAttr := slog.Any("defRef", snap.TypeRef)
	s.log.Debug("modification started", refAttr)
	rec, err := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DefRec, error) {
		return s.typeDefRepo.GetRecByRef(txCtx, conn, snap.TypeRef)
	})
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	if snap.TypeRef.DescRN != rec.TypeRef.DescRN {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, errConcurrentModification(snap.TypeRef.DescRN, rec.TypeRef.DescRN)
	}
	curSnap, err := s.retrieveSnap(ctx, rec)
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	if CheckSpec(snap.DefSpec.TypeExp, curSnap.DefSpec.TypeExp) == nil {
		return curSnap, nil
	}
	newExp, err := ConvertSpecToRec(snap.DefSpec.TypeExp)
	if err != nil {
		return DefSnap{}, err
	}
	expectedRoot := descroot.SemRef{DescID: rec.TypeRef.DescID, DescRN: snap.TypeRef.DescRN}
	snap.TypeRef.DescRN = seqnum.Next(snap.TypeRef.DescRN)
	rec.ExpCK = newExp.Key()
	err = s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if err := s.descRootRepo.TouchRef(txCtx, conn, expectedRoot); err != nil {
			return err
		}
		if err := s.typeExpRepo.AddRec(txCtx, conn, newExp); err != nil {
			return err
		}
		return s.typeDefRepo.ModifyRec(txCtx, conn, rec)
	})
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	s.log.Debug("modification succeed", refAttr)
	return snap, nil
}

func (s *service) RetrieveSnap(ctx context.Context, ref typesem.SemRef) (_ DefSnap, err error) {
	rec, getErr := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DefRec, error) {
		return s.typeDefRepo.GetRecByRef(txCtx, conn, ref)
	})
	if getErr != nil {
		s.log.Error("retrieval failed", slog.Any("ref", ref))
		return DefSnap{}, getErr
	}
	return s.retrieveSnap(ctx, rec)
}

func (s *service) retrieveSnap(ctx context.Context, rec DefRec) (_ DefSnap, err error) { // revive:disable-line:confusing-naming
	expRec, err := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (ExpRec, error) {
		return s.typeExpRepo.SelectRecByCK(txCtx, conn, rec.ExpCK)
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("ref", rec.TypeRef))
		return DefSnap{}, err
	}
	return DefSnap{
		TypeRef: rec.TypeRef,
		DefSpec: DefSpec{TypeExp: ConvertRecToSpec(expRec)},
	}, nil
}

func (s *service) RetreiveRefs(ctx context.Context) (refs []typesem.SemRef, err error) {
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

func ErrPolarityUnexpected(got ExpRec) error {
	return fmt.Errorf("root polarity unexpected: %v", got.Pol())
}

func ErrPolarityMismatch(a, b ExpRec) error {
	return fmt.Errorf("root polarity mismatch: %v != %v", a.Pol(), b.Pol())
}
