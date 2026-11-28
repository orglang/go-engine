package core

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/progroot"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

type API interface {
	Create(context.Context, DefSpec) (DefSnap, error)
	Modify(context.Context, DefSnap) (DefSnap, error)
	RetrieveSnap(context.Context, identity.ADT) (DefSnap, error)
	RetreiveRefs(context.Context) ([]identity.ADT, error)
}

type DefSpec struct {
	TypeQN  uniqsym.ADT
	TypeExpr ExpSpec
}

type DefRec struct {
	TypeQN  uniqsym.ADT
	TypeID  identity.ADT
	ProgRef progroot.RootRef
	ExpCK   contkey.ADT
}

type DefSnap struct {
	TypeID  identity.ADT
	ProgRef progroot.RootRef
	DefSpec DefSpec
}

type Context struct {
	Assets map[symbol.ADT]ExpRec
	Liabs  map[symbol.ADT]ExpRec
}

type service struct {
	typeDefRepo  Repo
	progRootRepo progroot.Repo
	typeExprrRepo  TypeExprRepo
	transactor   tx.Transactor
	log          *slog.Logger
}

func newAPI() API { return new(service) }

func NewService(typeDefRepo Repo, progRootRepo progroot.Repo, typeExprrRepo TypeExprRepo, transactor tx.Transactor, log *slog.Logger) API {
	return &service{typeDefRepo: typeDefRepo, progRootRepo: progRootRepo, typeExprrRepo: typeExprrRepo, transactor: transactor, log: log}
}

func (s *service) Create(ctx context.Context, spec DefSpec) (_ DefSnap, err error) {
	qnAttr := slog.Any("qn", spec.TypeQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newExp, err := ConvertSpecToRec(spec.TypeExpr)
	if err != nil {
		return DefSnap{}, err
	}
	newDef := DefRec{TypeID: identity.New(), ProgRef: progroot.New(), TypeQN: spec.TypeQN, ExpCK: newExp.Key()}
	err = s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if err := s.progRootRepo.AddRef(txCtx, conn, newDef.ProgRef); err != nil {
			return err
		}
		if err := s.typeExprrRepo.AddRec(txCtx, conn, newExp); err != nil {
			return err
		}
		return s.typeDefRepo.AddRec(txCtx, conn, newDef)
	})
	if err != nil {
		s.log.Error("creation failed", qnAttr)
		return DefSnap{}, err
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("typeID", newDef.TypeID), slog.Any("progRef", newDef.ProgRef))
	return DefSnap{TypeID: newDef.TypeID, ProgRef: newDef.ProgRef, DefSpec: spec}, nil
}

func (s *service) Modify(ctx context.Context, snap DefSnap) (_ DefSnap, err error) {
	refAttr := slog.Any("typeID", snap.TypeID)
	s.log.Debug("modification started", refAttr)
	rec, err := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DefRec, error) {
		return s.typeDefRepo.GetRecByID(txCtx, conn, snap.TypeID)
	})
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	if snap.ProgRef.ProgRN != rec.ProgRef.ProgRN {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, errConcurrentModification(snap.ProgRef.ProgRN, rec.ProgRef.ProgRN)
	}
	curSnap, err := s.snapFromRec(ctx, rec)
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	if CheckSpec(snap.DefSpec.TypeExpr, curSnap.DefSpec.TypeExpr) == nil {
		return curSnap, nil
	}
	newExp, err := ConvertSpecToRec(snap.DefSpec.TypeExpr)
	if err != nil {
		return DefSnap{}, err
	}
	expectedRoot := progroot.RootRef{ProgID: rec.ProgRef.ProgID, ProgRN: snap.ProgRef.ProgRN}
	snap.ProgRef.ProgRN = seqnum.Next(snap.ProgRef.ProgRN)
	rec.ExpCK = newExp.Key()
	err = s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if err := s.progRootRepo.TouchRef(txCtx, conn, expectedRoot); err != nil {
			return err
		}
		if err := s.typeExprrRepo.AddRec(txCtx, conn, newExp); err != nil {
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

func (s *service) RetrieveSnap(ctx context.Context, id identity.ADT) (_ DefSnap, err error) {
	rec, getErr := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (DefRec, error) {
		return s.typeDefRepo.GetRecByID(txCtx, conn, id)
	})
	if getErr != nil {
		s.log.Error("retrieval failed", slog.Any("typeID", id))
		return DefSnap{}, getErr
	}
	return s.snapFromRec(ctx, rec)
}

func (s *service) snapFromRec(ctx context.Context, rec DefRec) (_ DefSnap, err error) {
	expRec, err := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (ExpRec, error) {
		return s.typeExprrRepo.SelectRecByCK(txCtx, conn, rec.ExpCK)
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("typeID", rec.TypeID))
		return DefSnap{}, err
	}
	return DefSnap{TypeID: rec.TypeID, ProgRef: rec.ProgRef, DefSpec: DefSpec{TypeQN: rec.TypeQN, TypeExpr: ConvertRecToSpec(expRec)}}, nil
}

func (s *service) RetreiveRefs(ctx context.Context) (refs []identity.ADT, err error) {
	refs, err = tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) ([]identity.ADT, error) {
		return s.typeDefRepo.GetIDs(txCtx, conn)
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

func ErrSymMissingInEnv(want uniqsym.ADT) error { return fmt.Errorf("root missing in env: %v", want) }
func errConcurrentModification(got, want seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: want revision %v, got revision %v", want, got)
}
func ErrOptimisticUpdate(got seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
func ErrDoesNotExist(want identity.ADT) error { return fmt.Errorf("root doesn't exist: %v", want) }
func ErrMissingInEnv(want contkey.ADT) error  { return fmt.Errorf("exp missing in env: %v", want) }
func ErrMissingInCfg(want identity.ADT) error { return fmt.Errorf("root missing in cfg: %v", want) }
func ErrMissingInCtx(want symbol.ADT) error   { return fmt.Errorf("root missing in ctx: %v", want) }
func ErrPolarityUnexpected(got ExpRec) error {
	return fmt.Errorf("root polarity unexpected: %v", got.Pol())
}
func ErrPolarityMismatch(a, b ExpRec) error {
	return fmt.Errorf("root polarity mismatch: %v != %v", a.Pol(), b.Pol())
}
