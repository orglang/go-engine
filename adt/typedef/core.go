package typedef

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/synonym"
	"orglang/go-engine/adt/typeexp"
	"orglang/go-engine/adt/uniqref"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

type API interface {
	Incept(uniqsym.ADT) (DefRef, error)
	Create(DefSpec) (DefSnap, error)
	Modify(DefSnap) (DefSnap, error)
	RetrieveSnap(DefRef) (DefSnap, error)
	retrieveSnap(DefRec) (DefSnap, error)
	RetreiveRefs() ([]DefRef, error)
}

type DefRef = uniqref.ADT

type DefSpec struct {
	TypeQN uniqsym.ADT
	TypeES typeexp.ExpSpec
}

// aka TpDef
type DefRec struct {
	DefRef DefRef
	SynVK  valkey.ADT
	ExpVK  valkey.ADT
}

type DefSnap struct {
	DefRef  DefRef
	DefSpec DefSpec
}

type Context struct {
	Assets map[symbol.ADT]typeexp.ExpRec
	Liabs  map[symbol.ADT]typeexp.ExpRec
}

type service struct {
	typeDefs Repo
	typeExps typeexp.Repo
	synonyms synonym.Repo
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
	synDecs synonym.Repo,
	operator db.Operator,
	l *slog.Logger,
) *service {
	return &service{typeDefs, typeExps, synDecs, operator, l}
}

func (s *service) Incept(typeQN uniqsym.ADT) (_ DefRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("typeQN", typeQN)
	s.log.Debug("inception started", qnAttr)
	synVK, err := typeQN.Key()
	if err != nil {
		return DefRef{}, err
	}
	newSyn := synonym.Rec{SynQN: typeQN, SynVK: synVK}
	newDef := DefRec{DefRef: uniqref.New(), SynVK: synVK}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.synonyms.InsertRec(ds, newSyn)
		if err != nil {
			return err
		}
		err = s.typeDefs.InsertRec(ds, newDef)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("inception failed", qnAttr)
		return DefRef{}, err
	}
	s.log.Debug("inception succeed", qnAttr, slog.Any("defRef", newDef.DefRef))
	return newDef.DefRef, nil
}

func (s *service) Create(spec DefSpec) (_ DefSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("typeQN", spec.TypeQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	synVK, err := spec.TypeQN.Key()
	if err != nil {
		return DefSnap{}, err
	}
	newSyn := synonym.Rec{SynQN: spec.TypeQN, SynVK: synVK}
	newExp, err := typeexp.ConvertSpecToRec(spec.TypeES)
	if err != nil {
		return DefSnap{}, err
	}
	newDef := DefRec{
		DefRef: uniqref.New(),
		SynVK:  synVK,
		ExpVK:  newExp.Key(),
	}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.synonyms.InsertRec(ds, newSyn)
		if err != nil {
			return err
		}
		err = s.typeExps.InsertRec(ds, newExp, newDef.DefRef)
		if err != nil {
			return err
		}
		err = s.typeDefs.InsertRec(ds, newDef)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("creation failed", qnAttr)
		return DefSnap{}, err
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("defRef", newDef.DefRef))
	return DefSnap{
		DefRef:  newDef.DefRef,
		DefSpec: spec,
	}, nil
}

func (s *service) Modify(snap DefSnap) (_ DefSnap, err error) {
	ctx := context.Background()
	refAttr := slog.Any("defRef", snap.DefRef)
	s.log.Debug("modification started", refAttr)
	var rec DefRec
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		rec, err = s.typeDefs.SelectRecByRef(ds, snap.DefRef)
		return err
	})
	if err != nil {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, err
	}
	if snap.DefRef.RN != rec.DefRef.RN {
		s.log.Error("modification failed", refAttr)
		return DefSnap{}, errConcurrentModification(snap.DefRef.RN, rec.DefRef.RN)
	}
	snap.DefRef.RN = revnum.Next(snap.DefRef.RN)
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
			err = s.typeExps.InsertRec(ds, newExp, snap.DefRef)
			if err != nil {
				return err
			}
			rec.ExpVK = newExp.Key()
			rec.DefRef.RN = snap.DefRef.RN
		}
		if rec.DefRef.RN == snap.DefRef.RN {
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

func (s *service) RetrieveSnap(ref DefRef) (_ DefSnap, err error) {
	ctx := context.Background()
	var rec DefRec
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		rec, err = s.typeDefs.SelectRecByRef(ds, ref)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("defID", ref))
		return DefSnap{}, err
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
		s.log.Error("retrieval failed", slog.Any("defRef", rec.DefRef))
		return DefSnap{}, err
	}
	return DefSnap{
		DefRef:  rec.DefRef,
		DefSpec: DefSpec{TypeES: typeexp.ConvertRecToSpec(expRec)},
	}, nil
}

func (s *service) RetreiveRefs() (refs []DefRef, err error) {
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
