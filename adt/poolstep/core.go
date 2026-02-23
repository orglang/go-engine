package poolstep

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/poolexec"
	"orglang/go-engine/adt/poolexp"
	"orglang/go-engine/adt/procdec"
	"orglang/go-engine/adt/procexec"
)

type API interface {
	Take(StepSpec) error
	Spawn(StepSpec) (implsem.SemRef, error)
}

type StepSpec struct {
	ImplRef implsem.SemRef
	PoolES  poolexp.ExpSpec
}

type service struct {
	implSems  implsem.Repo
	procExecs procexec.Repo
	poolExecs poolexec.Repo
	procDecs  procdec.Repo
	operator  db.Operator
	log       *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	implSems implsem.Repo,
	procExecs procexec.Repo,
	poolExecs poolexec.Repo,
	procDecs procdec.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{implSems, procExecs, poolExecs, procDecs, operator, log.With(name)}
}

func (s *service) Take(spec StepSpec) (err error) {
	refAttr := slog.Any("pool", spec.ImplRef)
	s.log.Debug("step taking started", refAttr, slog.Any("exp", spec.PoolES))
	// descEnv (xact + pool)
	// 1. xactDefs + xactExps
	// 2. poolDecs
	// implCtx (xact + pool)
	// 1. poolExecs

	// xactDefs
	// 1. отдельно от xactExps
	// 2. вместе в снепшоте

	// descEnv (type + proc)
	// 1. typeDefs + typeExps
	// 2. procDecs + procDefs
	// implCtx (type + proc)
	// 1. procExecs
	s.log.Debug("step taking succeed", refAttr)
	return nil
}

func (s *service) Spawn(spec StepSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	poolAttr := slog.Any("pool", spec.ImplRef)
	s.log.Debug("proc spawning started", poolAttr, slog.Any("exp", spec.PoolES))
	spawn, ok := spec.PoolES.(poolexp.SpawnSpec)
	if !ok {
		panic(poolexp.ErrExpTypeUnexpected(spec.PoolES))
	}
	var procDec procdec.DecSnap
	selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		procDec, err = s.procDecs.SelectSnap(ds, spawn.ProcDescRef)
		return err
	})
	if selectErr != nil {
		s.log.Error("proc spawning failed", poolAttr)
		return implsem.SemRef{}, selectErr
	}
	newRef := implsem.NewRef()
	newImpl := implsem.SemRec{Ref: newRef, Kind: implsem.Proc}
	newExec := procexec.ExecRec{ImplRef: newRef, ChnlPH: procDec.ProviderVR.ChnlPH}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSems.InsertRec(ds, newImpl)
		if err != nil {
			return err
		}
		err = s.procExecs.InsertRec(ds, newExec)
		if err != nil {
			return err
		}
		return s.poolExecs.TouchRec(ds, spec.ImplRef)
	})
	if transactErr != nil {
		s.log.Error("proc spawning failed", poolAttr)
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("proc spawning succeed", poolAttr, slog.Any("proc", newRef))
	return newRef, nil
}
