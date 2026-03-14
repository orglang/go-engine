package poolexec

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/option"
	"orglang/go-engine/adt/poolconn"
	"orglang/go-engine/adt/poolctx"
	"orglang/go-engine/adt/pooldec"
	"orglang/go-engine/adt/poolenv"
	"orglang/go-engine/adt/poolexp"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/poolvar"
	"orglang/go-engine/adt/procdec"
	"orglang/go-engine/adt/procdef"
	"orglang/go-engine/adt/procexec"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/typedef"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/xactdef"
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Run(ExecSpec) (implsem.SemRef, error) // aka Create
	Take(poolstep.StepSpec) error
	Spawn(poolstep.StepSpec) (implsem.SemRef, error)
	RetrieveSnap(implsem.SemRef) (ExecSnap, error)
}

type ExecSpec struct {
	// ссылка на декларацию пула
	DescQN uniqsym.ADT
	// внутренняя и внешняя ссылки на вновь создаваемый пул
	LiabVar implvar.VarSpec
	// внутренние и внешние ссылки на ранее созданные пулы
	AssetVars []implvar.VarSpec
}

type ExecRec struct {
	ImplRef implsem.SemRef
}

type ExecSnap struct {
	ImplRef    implsem.SemRef
	StructVars map[symbol.ADT]implvar.VarRec
	LinearVars map[symbol.ADT]implvar.VarRec
}

type ExecMod struct {
	ImplRef    implsem.SemRef
	LinearVars []implvar.VarRec
}

type service struct {
	poolExecs Repo
	implSems  implsem.Repo
	commSems  commsem.Repo
	poolDecs  pooldec.Repo
	poolVars  poolvar.Repo
	poolConns poolconn.Repo
	poolSteps poolstep.Repo
	xactDefs  xactexp.Repo
	xactExps  xactexp.Repo
	procDecs  procdec.Repo
	procExecs procexec.Repo
	operator  db.Operator
	log       *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	poolExecs Repo,
	implSems implsem.Repo,
	commSems commsem.Repo,
	poolDecs pooldec.Repo,
	poolVars poolvar.Repo,
	poolConns poolconn.Repo,
	poolSteps poolstep.Repo,
	xactDefs xactexp.Repo,
	xactExps xactexp.Repo,
	procDecs procdec.Repo,
	procExecs procexec.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{poolExecs, implSems, commSems, poolDecs, poolVars, poolConns, poolSteps, xactDefs, xactExps, procDecs, procExecs, operator, log.With(name)}
}

func (s *service) Run(spec ExecSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	vsAttr := slog.Any("varSpec", spec.LiabVar)
	s.log.Debug("starting creation...", vsAttr, slog.Any("expSpec", spec))
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars))
	for _, assetVar := range spec.AssetVars {
		if assetVar.ImplQN == spec.LiabVar.ImplQN {
			continue
		}
		assetQNs = append(assetQNs, assetVar.ImplQN)
	}
	var liabVars map[uniqsym.ADT]commsem.SemRef
	selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		liabVars, err = s.poolConns.SelectRefsByQNs(ds, assetQNs)
		return err
	})
	if selectErr != nil {
		s.log.Error("creation failed", vsAttr)
		return implsem.SemRef{}, selectErr
	}
	implRef := implsem.NewRef()
	newImpl := implsem.SemRec{ImplRef: implRef, ImplQN: spec.LiabVar.ImplQN, Kind: implsem.Pool}
	newExec := ExecRec{ImplRef: implRef}
	newExch := commsem.SemRec{CommRef: commsem.NewRef(), Kind: commsem.Pool}
	newConn := poolconn.ConnRec{CommRef: newExch.CommRef}
	newLiab := implvar.VarRec{
		ImplRef: implRef,
		CommRef: newExch.CommRef,
		ChnlID:  identity.New(),
		ChnlPH:  spec.LiabVar.ChnlPH,
		ChnlBS:  implvar.Provider,
		// TODO: заполнить ExpVK
	}
	newAssets := make([]implvar.VarRec, 0, len(spec.AssetVars)+1)
	for _, assetVar := range spec.AssetVars {
		var commRef commsem.SemRef
		if assetVar.ImplQN == spec.LiabVar.ImplQN {
			commRef = newExch.CommRef
		} else {
			commRef = liabVars[assetVar.ImplQN]
		}
		newAssets = append(newAssets, implvar.VarRec{
			ImplRef: implRef,
			CommRef: commRef,
			ChnlID:  newLiab.ChnlID,
			ChnlPH:  assetVar.ChnlPH,
			ChnlBS:  implvar.Client,
			// TODO: заполнить ExpVK
		})
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSems.InsertRec(ds, newImpl)
		if err != nil {
			return err
		}
		err = s.poolExecs.InsertRec(ds, newExec)
		if err != nil {
			return err
		}
		err = s.poolVars.InsertRecs(ds, append(newAssets, newLiab))
		if err != nil {
			return err
		}
		err = s.commSems.InsertRec(ds, newExch)
		if err != nil {
			return err
		}
		return s.poolConns.InsertRec(ds, newConn)
	})
	if transactErr != nil {
		s.log.Error("creation failed", vsAttr)
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", vsAttr, slog.Any("ref", implRef))
	return implRef, nil
}

func (s *service) RetrieveSnap(ref implsem.SemRef) (snap ExecSnap, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		snap, err = s.poolExecs.SelectSubs(ds, ref)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("ref", ref))
		return ExecSnap{}, err
	}
	return snap, nil
}

func (s *service) Spawn(spec poolstep.StepSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	poolAttr := slog.Any("pool", spec.ImplRef)
	s.log.Debug("proc spawning started", poolAttr, slog.Any("exp", spec.PoolExp))
	spawn, ok := spec.PoolExp.(poolexp.SpawnSpec2)
	if !ok {
		panic(poolexp.ErrSpecTypeUnexpected(spec.PoolExp))
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
	newImpl := implsem.SemRec{ImplRef: newRef, Kind: implsem.Proc}
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

func (s *service) Take(spec poolstep.StepSpec) (err error) {
	ctx := context.Background()
	implAttr := slog.Any("impl", spec.ImplRef)
	s.log.Debug("step taking started", implAttr, slog.Any("exp", spec.PoolExp))
	var envSnap poolenv.EnvSnap
	envSpec := poolenv.ConvertExpToEnv(spec.PoolExp)
	selectErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
		envSnap, err = s.poolSteps.SelectEnvSnapByEnvSpec(ds, envSpec)
		return err
	})
	if selectErr1 != nil {
		return selectErr1
	}
	var ctxSnap poolctx.CtxSnap
	ctxSpec := poolctx.ConvertExpToSpec()
	selectErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
		ctxSnap, err = s.poolSteps.SelectCtxSnapByCtxSpec(ds, ctxSpec)
		return err
	})
	if selectErr2 != nil {
		return selectErr2
	}
	var execSnap ExecSnap
	selectErr3 := s.operator.Implicit(ctx, func(ds db.Source) error {
		execSnap, err = s.poolExecs.SelectSnap(ds, spec.ImplRef)
		return err
	})
	if selectErr3 != nil {
		return selectErr3
	}
	execMod, connMod, takeErr := s.take(envSnap, ctxSnap, execSnap, spec.PoolExp)
	if takeErr != nil {
		return takeErr
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.poolSteps.InsertRecs(ds, connMod.Steps)
		if err != nil {
			return err
		}
		err = s.poolConns.UpdateRec(ds, connMod)
		if err != nil {
			return err
		}
		err = s.commSems.TouchRec(ds, connMod.CommRef)
		if err != nil {
			return err
		}
		err = s.poolVars.InsertRecs(ds, execMod.LinearVars)
		if err != nil {
			return err
		}
		return s.implSems.TouchRec(ds, execMod.ImplRef)
	})
	if transactErr != nil {
		return transactErr
	}
	s.log.Debug("step taking succeed", implAttr)
	return nil
}

func (s *service) take(
	envSnap poolenv.EnvSnap,
	ctxSnap poolctx.CtxSnap,
	execSnap ExecSnap,
	es poolexp.ExpSpec,
) (
	execMod ExecMod,
	connMod poolconn.ConnMod,
	err error,
) {
	ctx := context.Background()
	execMod.ImplRef = execSnap.ImplRef
	implAttr := slog.Any("impl", execSnap.ImplRef)
	switch expSpec := es.(type) {
	case poolexp.AcceptSpec:
		commChnl, ok := execSnap.StructVars[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implAttr)
			return execMod, connMod, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		commAttr := slog.Any("comm", commChnl.CommRef)
		// вычисляем следующее состояние
		xactExp, ok := envSnap.XactExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", implAttr, commAttr)
			return execMod, connMod, xactdef.ErrMissingInEnv(commChnl.ExpVK)
		}
		nextExpVK := xactExp.(xactexp.ProdRec).Next()
		// получаем снепшот соединения
		var connSnap poolconn.ConnSnap
		selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.poolConns.SelectSnapByQry(ds, poolconn.ConnQuery{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if selectErr != nil {
			s.log.Error("step taking failed", implAttr, commAttr)
			return execMod, connMod, selectErr
		}
		connMod.CommRef = connSnap.CommRef
		subscription := connSnap.NextStep()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступодателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.VarRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем сообщение доступодателя
			connMod.Steps = append(connMod.Steps, poolstep.PubRec{
				CommRef: connSnap.CommRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: poolexp.AcceptRec{
					ContChnlID: newChnlID,
				},
			})
			s.log.Debug("taking half done", implAttr, commAttr)
			return execMod, connMod, nil
		}
		acquisition, ok := subscription.(poolstep.SubRec)
		if !ok {
			panic(poolstep.ErrRecTypeUnexpected(subscription))
		}
		switch contExp := acquisition.ContExp.(type) {
		case poolexp.AcquireRec:
			// сдвигаем офсет коммуникации
			connMod.CommON = option.Some(acquisition.CommRef.CommRN)
			// вяжем продолжение доступодателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.VarRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			s.log.Debug("step taking succeed", implAttr, commAttr)
			return execMod, connMod, nil
		default:
			panic(poolexp.ErrRecTypeUnexpected(acquisition.ContExp))
		}
	case poolexp.AcquireSpec:
		commChnl, ok := execSnap.StructVars[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implAttr)
			return execMod, connMod, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		commAttr := slog.Any("comm", commChnl.CommRef)
		xactExp, ok := envSnap.XactExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", implAttr, commAttr)
			return execMod, connMod, typedef.ErrMissingInEnv(commChnl.ExpVK)
		}
		nextExpVK := xactExp.(xactexp.ProdRec).Next()
		// получаем снепшот соединения
		var connSnap poolconn.ConnSnap
		selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.poolConns.SelectSnapByQry(ds, poolconn.ConnQuery{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if selectErr != nil {
			s.log.Error("step taking failed", implAttr, commAttr)
			return execMod, connMod, selectErr
		}
		connMod.CommRef = connSnap.CommRef
		publication := connSnap.NextStep()
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступополучателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.VarRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем подписку доступополучателя
			connMod.Steps = append(connMod.Steps, poolstep.SubRec{
				CommRef: connSnap.CommRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: poolexp.AcceptRec{
					ContChnlID: newChnlID,
				},
			})
			s.log.Debug("taking half done", implAttr, commAttr)
			return execMod, connMod, nil
		}
		acception, ok := publication.(poolstep.PubRec)
		if !ok {
			panic(poolstep.ErrRecTypeUnexpected(publication))
		}
		switch valExp := acception.ValExp.(type) {
		case poolexp.AcceptRec:
			// сдвигаем офсет коммуникации
			connMod.CommON = option.Some(acception.CommRef.CommRN)
			// вяжем продолжение доступополучателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.VarRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			s.log.Debug("step taking succeed", implAttr, commAttr)
			return execMod, connMod, nil
		default:
			panic(poolexp.ErrRecTypeUnexpected(acception.ValExp))
		}
	default:
		panic(poolexp.ErrSpecTypeUnexpected(expSpec))
	}
}
