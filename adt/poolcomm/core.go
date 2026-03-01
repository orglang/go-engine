package poolcomm

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/poolcfg"
	"orglang/go-engine/adt/poolctx"
	"orglang/go-engine/adt/poolenv"
	"orglang/go-engine/adt/poolexec"
	"orglang/go-engine/adt/poolexp"
	"orglang/go-engine/adt/poolvar"
	"orglang/go-engine/adt/procdec"
	"orglang/go-engine/adt/procdef"
	"orglang/go-engine/adt/procexec"
	"orglang/go-engine/adt/typedef"
	"orglang/go-engine/adt/xactdef"
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Take(CommSpec) error
	Spawn(CommSpec) (implsem.SemRef, error)
}

type CommSpec struct {
	ImplRef implsem.SemRef
	PoolES  poolexp.ExpSpec
}

type CommMod struct {
	Refs  []implsem.SemRef
	Vars  []implvar.VarRec
	Comms []CommRec
}

type CommSnap struct {
	Vars  []implvar.VarRec
	Comms []CommRec
}

// communication (aka Sem)
type CommRec interface {
	comm()
}

// publication (aka Msg)
type PubRec struct {
	ImplRef implsem.SemRef
	ChnlID  identity.ADT
	CommRef commsem.SemRef
	ValExp  poolexp.ExpRec
}

func (r PubRec) comm() {}

// subscription (aka Service)
type SubRec struct {
	ImplRef implsem.SemRef
	ChnlID  identity.ADT
	CommRef commsem.SemRef
	ContExp poolexp.ExpRec
}

func (r SubRec) comm() {}

type service struct {
	poolComms Repo
	implSems  implsem.Repo
	procExecs procexec.Repo
	poolExecs poolexec.Repo
	poolVars  poolvar.Repo
	procDecs  procdec.Repo
	operator  db.Operator
	log       *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	poolComms Repo,
	implSems implsem.Repo,
	procExecs procexec.Repo,
	poolExecs poolexec.Repo,
	poolVars poolvar.Repo,
	procDecs procdec.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{poolComms, implSems, procExecs, poolExecs, poolVars, procDecs, operator, log.With(name)}
}

func (s *service) Take(spec CommSpec) (err error) {
	ctx := context.Background()
	refAttr := slog.Any("pool", spec.ImplRef)
	s.log.Debug("step taking started", refAttr, slog.Any("exp", spec.PoolES))
	var envSnap poolenv.EnvSnap
	envSpec := poolenv.ConvertExpToEnv(spec.PoolES)
	selectErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
		envSnap, err = s.poolComms.SelectEnvSnapByEnvSpec(ds, envSpec)
		return err
	})
	if selectErr1 != nil {
		return selectErr1
	}
	var ctxSnap poolctx.CtxSnap
	ctxSpec := poolctx.ConvertExpToSpec()
	selectErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
		ctxSnap, err = s.poolComms.SelectCtxSnapByCtxSpec(ds, ctxSpec)
		return err
	})
	if selectErr2 != nil {
		return selectErr2
	}
	var cfgSnap poolcfg.CfgSnap
	cfgSpec := poolcfg.ConvertExpToSpec()
	selectErr3 := s.operator.Implicit(ctx, func(ds db.Source) error {
		cfgSnap, err = s.poolComms.SelectCfgSnapBySpec(ds, cfgSpec)
		return err
	})
	if selectErr3 != nil {
		return selectErr3
	}
	mod, takeErr := s.take(envSnap, ctxSnap, cfgSnap, spec)
	if takeErr != nil {
		return takeErr
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.poolVars.InsertRecs(ds, mod.Vars)
		if err != nil {
			return err
		}
		err = s.poolComms.InsertRecs(ds, mod.Comms)
		if err != nil {
			return err
		}
		return s.poolExecs.TouchRecs(ds, mod.Refs)
	})
	if transactErr != nil {
		return transactErr
	}
	s.log.Debug("step taking succeed", refAttr)
	return nil
}

func (s *service) Spawn(spec CommSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	poolAttr := slog.Any("pool", spec.ImplRef)
	s.log.Debug("proc spawning started", poolAttr, slog.Any("exp", spec.PoolES))
	spawn, ok := spec.PoolES.(poolexp.SpawnSpec2)
	if !ok {
		panic(poolexp.ErrSpecTypeUnexpected(spec.PoolES))
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

func (s *service) take(
	envSnap poolenv.EnvSnap,
	ctxSnap poolctx.CtxSnap,
	cfgSnap poolcfg.CfgSnap,
	commSpec CommSpec,
) (
	commMod CommMod,
	err error,
) {
	ctx := context.Background()
	// на чем синхронизируемся при коммуникации? на исполнении или на канале?
	// что записываем в качестве офсета? ревизию исполнения или ревизию канала?
	commMod.Refs = append(commMod.Refs, commSpec.ImplRef)
	switch expSpec := commSpec.PoolES.(type) {
	case poolexp.AcceptSpec:
		commChnl, ok := cfgSnap.SharedVars[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", slog.Any("pool", commSpec.ImplRef))
			return CommMod{}, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		refAttr := slog.Any("chnl", commChnl.ChnlRef)
		var subscription CommRec
		selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			subscription, err = s.poolComms.SelectRecByRef(ds, commChnl.ChnlRef)
			return err
		})
		if selectErr != nil {
			return CommMod{}, selectErr
		}
		// сдвигаем офсет доступодателя
		commMod.Vars = append(commMod.Vars, commChnl.Rewind(commSpec.ImplRef.ImplRN))
		// вычисляем следующее состояние
		xactExp, ok := envSnap.XactExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", refAttr)
			return CommMod{}, xactdef.ErrMissingInEnv(commChnl.ExpVK)
		}
		nextExpVK := xactExp.(xactexp.ProdRec).Next()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступодателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: commSpec.ImplRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  implvar.Provider,
				ExpVK:   nextExpVK,
			})
			// регистрируем сообщение доступодателя
			commMod.Comms = append(commMod.Comms, PubRec{
				ImplRef: commSpec.ImplRef,
				ChnlID:  commChnl.ChnlRef.ChnlID,
				ValExp: poolexp.AcceptRec{
					ContChnlID: newChnlID,
				},
			})
			s.log.Debug("taking half done", refAttr)
			return commMod, nil
		}
		acquisition, ok := subscription.(SubRec)
		if !ok {
			panic(ErrRecTypeUnexpected(subscription))
		}
		switch contExp := acquisition.ContExp.(type) {
		case poolexp.AcquireRec:
			// вяжем продолжение доступодателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: commSpec.ImplRef,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  implvar.Provider,
				ExpVK:   nextExpVK,
			})
			s.log.Debug("step taking succeed", refAttr)
			return commMod, nil
		default:
			panic(poolexp.ErrRecTypeUnexpected(acquisition.ContExp))
		}
	case poolexp.AcquireSpec:
		commChnl, ok := cfgSnap.SharedVars[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed")
			return CommMod{}, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		refAttr := slog.Any("chnl", commChnl.ChnlRef)
		var publication CommRec
		selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			publication, err = s.poolComms.SelectRecByRef(ds, commChnl.ChnlRef)
			return err
		})
		if selectErr != nil {
			return CommMod{}, selectErr
		}
		xactExp, ok := envSnap.XactExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", refAttr)
			return CommMod{}, typedef.ErrMissingInEnv(commChnl.ExpVK)
		}
		nextExpVK := xactExp.(xactexp.ProdRec).Next()
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступополучателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: commSpec.ImplRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  implvar.Client,
				ExpVK:   nextExpVK,
			})
			// регистрируем подписку доступополучателя
			commMod.Comms = append(commMod.Comms, SubRec{
				// ImplRef - это ссылка на исполнение доступополучателя
				// commChnl - это клиентская сторона канала
				// 1. если тут rn доступодателя, то он по идее устарел
				// 2. если тут rn доступополучателя, то он не подходит
				CommRef: commChnl.ChnlRef,
				ChnlID:  commChnl.ChnlRef.ChnlID,
				ContExp: poolexp.AcceptRec{
					ContChnlID: newChnlID,
				},
			})
			s.log.Debug("taking half done", refAttr)
			return commMod, nil
		}
		acception, ok := publication.(PubRec)
		if !ok {
			panic(ErrRecTypeUnexpected(publication))
		}
		switch valExp := acception.ValExp.(type) {
		case poolexp.AcceptRec:
			// вяжем продолжение доступополучателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: commSpec.ImplRef,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  implvar.Client,
				ExpVK:   nextExpVK,
			})
			s.log.Debug("step taking succeed", refAttr)
			return commMod, nil
		default:
			panic(poolexp.ErrRecTypeUnexpected(acception.ValExp))
		}
	default:
		panic(poolexp.ErrSpecTypeUnexpected(expSpec))
	}
}

func ErrRecTypeUnexpected(got CommRec) error {
	return fmt.Errorf("comm rec unexpected: %T", got)
}

func ErrRecTypeMismatch(got, want CommRec) error {
	return fmt.Errorf("comm rec mismatch: want %T, got %T", want, got)
}
