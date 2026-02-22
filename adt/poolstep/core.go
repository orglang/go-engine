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
	PoolIR implsem.SemRef
	PoolES poolexp.ExpSpec
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
	refAttr := slog.Any("ref", spec.PoolIR)
	s.log.Debug("starting taking...", refAttr, slog.Any("spec", spec.PoolES))
	s.log.Debug("taking succeed", refAttr)
	return nil
}

func (s *service) Spawn(spec StepSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	refAttr := slog.Any("ref", spec.PoolIR)
	s.log.Debug("starting spawning...", refAttr, slog.Any("spec", spec.PoolES))
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
		s.log.Error("creation failed", refAttr)
		return implsem.SemRef{}, selectErr
	}
	newRef := implsem.NewRef()
	newImpl := implsem.SemRec{Ref: newRef, Kind: implsem.Proc}
	newExec := procexec.ExecRec{ProcIR: newRef, ProviderPH: procDec.ProviderVR.ChnlPH}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSems.InsertRec(ds, newImpl)
		if err != nil {
			return err
		}
		err = s.procExecs.InsertRec(ds, newExec)
		if err != nil {
			return err
		}
		return s.poolExecs.TouchRec(ds, spec.PoolIR)
	})
	if transactErr != nil {
		s.log.Error("creation failed", refAttr)
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("spawning succeed", refAttr)
	return implsem.SemRef{}, nil
}
