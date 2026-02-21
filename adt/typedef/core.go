package typedef

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/typeexp"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

type API interface {
	Incept(uniqsym.ADT) (descsem.SemRef, error)
	Create(DefSpec) (DefSnap, error)
	Modify(DefSnap) (DefSnap, error)
	RetrieveSnap(descsem.SemRef) (DefSnap, error)
	retrieveSnap(DefRec) (DefSnap, error)
	RetreiveRefs() ([]descsem.SemRef, error)
}

type DefSpec struct {
	TypeQN uniqsym.ADT
	TypeES typeexp.ExpSpec
}

// aka TpDef
type DefRec struct {
	DescRef descsem.SemRef
	ExpVK   valkey.ADT
}

type DefSnap struct {
	DescRef descsem.SemRef
	DefSpec DefSpec
}

type Context struct {
	Assets map[symbol.ADT]typeexp.ExpRec
	Liabs  map[symbol.ADT]typeexp.ExpRec
}

type service struct {
	typeDefs Repo
	typeExps typeexp.Repo
	descSems descsem.Repo
	operator db.Operator
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	typeDefs Repo,
	typeExps typeexp.Repo,
	descSems descsem.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	return &service{typeDefs, typeExps, descSems, operator, log}
}

func (s *service) Incept(typeQN uniqsym.ADT) (_ descsem.SemRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", typeQN)
	s.log.Debug("starting inception...", qnAttr)
	newRef := descsem.NewRef()
	newBind := descsem.SemBind{DescQN: typeQN, DescID: newRef.DescID}
	newDesc := descsem.SemRec{Ref: descsem.NewRef(), Bind: newBind, Kind: descsem.Type}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descSems.InsertRec(ds, newDesc)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("inception failed", qnAttr)
		return descsem.SemRef{}, err
	}
	s.log.Debug("inception succeed", qnAttr, slog.Any("ref", newRef))
	return newRef, nil
}

func (s *service) Create(spec DefSpec) (_ DefSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", spec.TypeQN)
	s.log.Debug("starting creation...", qnAttr, slog.Any("spec", spec))
	newRef := descsem.NewRef()
	newBind := descsem.SemBind{DescQN: spec.TypeQN, DescID: newRef.DescID}
	newDesc := descsem.SemRec{Ref: descsem.NewRef(), Bind: newBind, Kind: descsem.Type}
	newExp, err := typeexp.ConvertSpecToRec(spec.TypeES)
	if err != nil {
		return DefSnap{}, err
	}
	newDef := DefRec{DescRef: newRef, ExpVK: newExp.Key()}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descSems.InsertRec(ds, newDesc)
		if err != nil {
			return err
		}
		err = s.typeExps.InsertRec(ds, newExp, newRef)
		if err != nil {
			return err
		}
		return s.typeDefs.InsertRec(ds, newDef)
	})
	if err != nil {
		s.log.Error("creation failed", qnAttr)
		return DefSnap{}, err
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newRef))
	return DefSnap{
		DescRef: newRef,
		DefSpec: spec,
	}, nil
}

func (s *service) Modify(snap DefSnap) (_ DefSnap, err error) {
	ctx := context.Background()
	refAttr := slog.Any("defRef", snap.DescRef)
	s.log.Debug("modification started", refAttr)
	var rec DefRec
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		rec, err = s.typeDefs.SelectRecByRef(ds, snap.DescRef)
		return err
	})
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	if snap.DescRef.DescRN != rec.DescRef.DescRN {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, errConcurrentModification(snap.DescRef.DescRN, rec.DescRef.DescRN)
	}
	snap.DescRef.DescRN = revnum.Next(snap.DescRef.DescRN)
	curSnap, err := s.retrieveSnap(rec)
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		if typeexp.CheckSpec(snap.DefSpec.TypeES, curSnap.DefSpec.TypeES) != nil {
			newExp, err := typeexp.ConvertSpecToRec(snap.DefSpec.TypeES)
			if err != nil {
				return err
			}
			err = s.typeExps.InsertRec(ds, newExp, snap.DescRef)
			if err != nil {
				return err
			}
			rec.ExpVK = newExp.Key()
			rec.DescRef.DescRN = snap.DescRef.DescRN
		}
		if rec.DescRef.DescRN == snap.DescRef.DescRN {
			err = s.typeDefs.Update(ds, rec)
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

func (s *service) RetrieveSnap(ref descsem.SemRef) (_ DefSnap, err error) {
	ctx := context.Background()
	var rec DefRec
	selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		rec, err = s.typeDefs.SelectRecByRef(ds, ref)
		return err
	})
	if selectErr != nil {
		s.log.Error("retrieval failed", slog.Any("ref", ref))
		return DefSnap{}, selectErr
	}
	return s.retrieveSnap(rec)
}

func (s *service) retrieveSnap(rec DefRec) (_ DefSnap, err error) { // revive:disable-line:confusing-naming
	ctx := context.Background()
	var expRec typeexp.ExpRec
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		expRec, err = s.typeExps.SelectRecByVK(ds, rec.ExpVK)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("ref", rec.DescRef))
		return DefSnap{}, err
	}
	return DefSnap{
		DescRef: rec.DescRef,
		DefSpec: DefSpec{TypeES: typeexp.ConvertRecToSpec(expRec)},
	}, nil
}

func (s *service) RetreiveRefs() (refs []descsem.SemRef, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		refs, err = s.typeDefs.SelectRefs(ds)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}

func CollectEnv(recs iter.Seq[DefRec]) []valkey.ADT {
	expIDs := []valkey.ADT{}
	for r := range recs {
		expIDs = append(expIDs, r.ExpVK)
	}
	return expIDs
}

func ErrSymMissingInEnv(want uniqsym.ADT) error {
	return fmt.Errorf("root missing in env: %v", want)
}

func errConcurrentModification(got revnum.ADT, want revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: want revision %v, got revision %v", want, got)
}

func errOptimisticUpdate(got revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}

func ErrDoesNotExist(want identity.ADT) error {
	return fmt.Errorf("root doesn't exist: %v", want)
}

func ErrMissingInEnv(want valkey.ADT) error {
	return fmt.Errorf("exp missing in env: %v", want)
}

func ErrMissingInCfg(want identity.ADT) error {
	return fmt.Errorf("root missing in cfg: %v", want)
}

func ErrMissingInCtx(want symbol.ADT) error {
	return fmt.Errorf("root missing in ctx: %v", want)
}
