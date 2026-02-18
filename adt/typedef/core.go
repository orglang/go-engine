package typedef

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descbind"
	"orglang/go-engine/adt/descexec"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/typeexp"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

type API interface {
	Incept(uniqsym.ADT) (descexec.ExecRef, error)
	Create(DefSpec) (DefSnap, error)
	Modify(DefSnap) (DefSnap, error)
	RetrieveSnap(descexec.ExecRef) (DefSnap, error)
	retrieveSnap(DefRec) (DefSnap, error)
	RetreiveRefs() ([]descexec.ExecRef, error)
}

type DefSpec struct {
	TypeQN uniqsym.ADT
	TypeES typeexp.ExpSpec
}

// aka TpDef
type DefRec struct {
	DescRef descexec.ExecRef
	ExpVK   valkey.ADT
}

type DefSnap struct {
	DescRef descexec.ExecRef
	DefSpec DefSpec
}

type Context struct {
	Assets map[symbol.ADT]typeexp.ExpRec
	Liabs  map[symbol.ADT]typeexp.ExpRec
}

type service struct {
	typeDefs  Repo
	typeExps  typeexp.Repo
	descExecs descexec.Repo
	descBinds descbind.Repo
	operator  db.Operator
	log       *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	typeDefs Repo,
	typeExps typeexp.Repo,
	descExecs descexec.Repo,
	descBinds descbind.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	return &service{typeDefs, typeExps, descExecs, descBinds, operator, log}
}

func (s *service) Incept(typeQN uniqsym.ADT) (_ descexec.ExecRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("typeQN", typeQN)
	s.log.Debug("inception started", qnAttr)
	newExec := descexec.ExecRec{Ref: descexec.NewExecRef(), Kind: descexec.Type}
	newBind := descbind.BindRec{DescQN: typeQN, DescID: newExec.Ref.DescID}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descBinds.InsertRec(ds, newBind)
		if err != nil {
			return err
		}
		err = s.descExecs.InsertRec(ds, newExec)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("inception failed", qnAttr)
		return descexec.ExecRef{}, err
	}
	s.log.Debug("inception succeed", qnAttr, slog.Any("ref", newExec.Ref))
	return newExec.Ref, nil
}

func (s *service) Create(spec DefSpec) (_ DefSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("typeQN", spec.TypeQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newExec := descexec.ExecRec{Ref: descexec.NewExecRef(), Kind: descexec.Type}
	newBind := descbind.BindRec{DescQN: spec.TypeQN, DescID: newExec.Ref.DescID}
	newExp, err := typeexp.ConvertSpecToRec(spec.TypeES)
	if err != nil {
		return DefSnap{}, err
	}
	newDef := DefRec{DescRef: newExec.Ref, ExpVK: newExp.Key()}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descBinds.InsertRec(ds, newBind)
		if err != nil {
			return err
		}
		err = s.descExecs.InsertRec(ds, newExec)
		if err != nil {
			return err
		}
		err = s.typeExps.InsertRec(ds, newExp, newExec.Ref)
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
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newExec.Ref))
	return DefSnap{
		DescRef: newExec.Ref,
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

func (s *service) RetrieveSnap(ref descexec.ExecRef) (_ DefSnap, err error) {
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

func (s *service) RetreiveRefs() (refs []descexec.ExecRef, err error) {
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
