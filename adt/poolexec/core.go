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
	"orglang/go-engine/adt/pooldec"
	"orglang/go-engine/adt/poolexp"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/poolvar"
	"orglang/go-engine/adt/procdec"
	"orglang/go-engine/adt/procdef"
	"orglang/go-engine/adt/procexec"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
	"orglang/go-engine/adt/xactdef"
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Run(ExecSpec) (implsem.SemRef, error) // aka Create
	Take(poolstep.StepSpec) error
	Spawn(poolstep.StepSpec) (implsem.SemRef, error)
	RetrieveSnap(implsem.SemRef) (ExecCtxSnap, error)
}

type ExecSpec struct {
	// ссылка на декларацию вновь создаваемого пула
	DescQN uniqsym.ADT
	// внутренняя и внешняя ссылки на вновь создаваемый пул
	LiabVar implvar.VarSpec
	// внутренние и внешние ссылки на ранее созданные пулы
	AssetVars []implvar.VarSpec
}

type ExecRec struct {
	ImplRef  implsem.SemRef
	LiabMode implvar.Mode
}

type ExecMod struct {
	ImplRef implsem.SemRef
	Vars    []implvar.VarRec
}

type ExecCtxSnap struct {
	ImplRef    implsem.SemRef
	StructVars map[symbol.ADT]implvar.StructRec
	LinearVars map[symbol.ADT]implvar.LinearRec
	XactExps   map[valkey.ADT]xactexp.ExpSpec
}

type ExecLiabSnap struct {
	ImplRef implsem.SemRef
	LiabVar implvar.VarRec
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
	specAttr := slog.Any("spec", spec)
	s.log.Debug("creation started", specAttr)
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars))
	for _, assetVar := range spec.AssetVars {
		if assetVar.ImplQN == spec.LiabVar.ImplQN {
			continue
		}
		assetQNs = append(assetQNs, assetVar.ImplQN)
	}
	var assetExecs map[uniqsym.ADT]ExecLiabSnap
	getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		assetExecs, err = s.poolExecs.GetSnapsByQNs(ds, assetQNs)
		return err
	})
	if getErr != nil {
		s.log.Error("creation failed", specAttr)
		return implsem.SemRef{}, getErr
	}
	newImplSem := implsem.SemRec{ImplRef: implsem.NewRef(), ImplQN: spec.LiabVar.ImplQN, Kind: implsem.Pool}
	newCommSem := commsem.SemRec{CommRef: commsem.NewRef(), Kind: commsem.Pool}
	newConn := poolconn.ConnRec{CommRef: newCommSem.CommRef}
	newExec := ExecRec{ImplRef: newImplSem.ImplRef, LiabMode: implvar.StructMode}
	newLiabVar := implvar.StructRec{
		ImplRef: newImplSem.ImplRef,
		CommRef: newCommSem.CommRef,
		ChnlID:  identity.New(),
		ChnlPH:  spec.LiabVar.ChnlPH,
		ChnlBS:  implvar.LiabSide,
		// TODO: заполнить ExpVK
	}
	newAssetVars := make([]implvar.VarRec, 0, len(spec.AssetVars)+1)
	for _, assetVar := range spec.AssetVars {
		var commRef commsem.SemRef
		var chnlID identity.ADT
		assetExec, ok := assetExecs[assetVar.ImplQN]
		if !ok && assetVar.ImplQN == spec.LiabVar.ImplQN {
			commRef = newCommSem.CommRef
			chnlID = newLiabVar.ChnlID
		} else {
			commRef = assetExec.LiabVar.GetCommRef()
			chnlID = assetExec.LiabVar.GetChnlID()
		}
		newAssetVars = append(newAssetVars, implvar.StructRec{
			ImplRef: newImplSem.ImplRef,
			CommRef: commRef,
			ChnlID:  chnlID,
			ChnlPH:  assetVar.ChnlPH,
			ChnlBS:  implvar.AssetSide,
			// TODO: заполнить ExpVK
		})
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSems.AddRec(ds, newImplSem)
		if err != nil {
			return err
		}
		err = s.poolExecs.AddRec(ds, newExec)
		if err != nil {
			return err
		}
		err = s.poolVars.AddRecs(ds, append(newAssetVars, newLiabVar))
		if err != nil {
			return err
		}
		err = s.commSems.AddRec(ds, newCommSem)
		if err != nil {
			return err
		}
		return s.poolConns.AddRec(ds, newConn)
	})
	if transactErr != nil {
		s.log.Error("creation failed", specAttr)
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", specAttr, slog.Any("ref", newImplSem.ImplRef))
	return newImplSem.ImplRef, nil
}

func (s *service) RetrieveSnap(ref implsem.SemRef) (snap ExecCtxSnap, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		snap, err = s.poolExecs.GetSnap(ds, ref)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("ref", ref))
		return ExecCtxSnap{}, err
	}
	return snap, nil
}

func (s *service) Spawn(spec poolstep.StepSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	implAttr := slog.Any("pool", spec.ImplRef)
	s.log.Debug("proc spawning started", implAttr, slog.Any("exp", spec.PoolExp))
	spawn, ok := spec.PoolExp.(poolexp.SpawnSpec2)
	if !ok {
		panic(poolexp.ErrSpecTypeUnexpected(spec.PoolExp))
	}
	var procDec procdec.DecSnap
	getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		procDec, err = s.procDecs.SelectSnap(ds, spawn.ProcDescRef)
		return err
	})
	if getErr != nil {
		s.log.Error("proc spawning failed", implAttr)
		return implsem.SemRef{}, getErr
	}
	newRef := implsem.NewRef()
	newImpl := implsem.SemRec{ImplRef: newRef, Kind: implsem.Proc}
	newExec := procexec.ExecRec{ImplRef: newRef, ChnlPH: procDec.LiabVar.ChnlPH}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSems.AddRec(ds, newImpl)
		if err != nil {
			return err
		}
		return s.procExecs.InsertRec(ds, newExec)
	})
	if transactErr != nil {
		s.log.Error("proc spawning failed", implAttr)
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("proc spawning succeed", implAttr, slog.Any("proc", newRef))
	return newRef, nil
}

func (s *service) Take(spec poolstep.StepSpec) (err error) {
	ctx := context.Background()
	implAttr := slog.Any("impl", spec.ImplRef)
	s.log.Debug("step taking started", implAttr, slog.Any("exp", spec.PoolExp))
	var execSnap ExecCtxSnap
	selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		execSnap, err = s.poolExecs.GetSnap(ds, spec.ImplRef)
		return err
	})
	if selectErr != nil {
		s.log.Error("step taking failed", implAttr)
		return selectErr
	}
	execMod, connMod, takeErr := s.take(execSnap, spec.PoolExp)
	if takeErr != nil {
		s.log.Error("step taking failed", implAttr)
		return takeErr
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.poolSteps.AddRecs(ds, connMod.Steps)
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
		err = s.poolVars.AddRecs(ds, execMod.Vars)
		if err != nil {
			return err
		}
		return s.implSems.TouchRec(ds, execMod.ImplRef)
	})
	if transactErr != nil {
		s.log.Error("step taking failed", implAttr)
		return transactErr
	}
	s.log.Debug("step taking succeed", implAttr)
	return nil
}

func (s *service) take(
	execSnap ExecCtxSnap,
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
		connMod.CommRef = commChnl.CommRef
		commAttr := slog.Any("comm", commChnl.CommRef)
		// вычисляем следующее состояние
		xactExp, ok := execSnap.XactExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", implAttr, commAttr)
			return execMod, connMod, xactdef.ErrMissingInEnv(commChnl.ExpVK)
		}
		nextExpVK := xactExp.(xactexp.ProdRec).Next()
		// получаем снепшот соединения
		var connSnap poolconn.ConnSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.poolConns.GetSnapByQry(ds, poolconn.ConnQry{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implAttr, commAttr)
			return execMod, connMod, getErr
		}
		subscription := connSnap.NextStep()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступодателя
			execMod.Vars = append(execMod.Vars, implvar.LinearRec{
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
			execMod.Vars = append(execMod.Vars, implvar.LinearRec{
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
		connMod.CommRef = commChnl.CommRef
		commAttr := slog.Any("comm", commChnl.CommRef)
		xactExp, ok := execSnap.XactExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", implAttr, commAttr)
			return execMod, connMod, xactdef.ErrMissingInEnv(commChnl.ExpVK)
		}
		nextExpVK := xactExp.(xactexp.ProdRec).Next()
		// получаем снепшот соединения
		var connSnap poolconn.ConnSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.poolConns.GetSnapByQry(ds, poolconn.ConnQry{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implAttr, commAttr)
			return execMod, connMod, getErr
		}
		connMod.CommRef = connSnap.CommRef
		publication := connSnap.NextStep()
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступополучателя
			execMod.Vars = append(execMod.Vars, implvar.LinearRec{
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
			execMod.Vars = append(execMod.Vars, implvar.LinearRec{
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
