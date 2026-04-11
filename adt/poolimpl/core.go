package poolimpl

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
	"orglang/go-engine/adt/poolcomm"
	"orglang/go-engine/adt/poolexec"
	"orglang/go-engine/adt/poolexp"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/poolvar"
	"orglang/go-engine/adt/procdef"
	"orglang/go-engine/adt/procexec"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Take(poolstep.StepSpec) error
	Spawn(poolstep.StepSpec) (implsem.SemRef, error)
}

type ImplRec struct {
	ImplRef    implsem.SemRef
	LiabMode   implvar.Mode
	StructVars []implvar.StructRec
	LinearVars []implvar.LinearRec
}

type ImplSnap struct {
	ImplRef    implsem.SemRef
	StructVars map[symbol.ADT]implvar.StructRec
	StructExps map[symbol.ADT]xactexp.ExpRec
	LinearVars map[symbol.ADT]implvar.LinearRec
	LinearExps map[symbol.ADT]xactexp.ExpRec
}

type ImplMod struct {
	Vars []implvar.VarRec
}

type ImplEff struct {
	Steps []poolstep.StepSpec
}

type service struct {
	poolImplRepo Repo
	poolImplExch Exch
	poolExecRepo poolexec.Repo
	implSemRepo  implsem.Repo
	commSemRepo  commsem.Repo
	poolVarRepo  poolvar.Repo
	poolCommRepo poolcomm.Repo
	poolStepRepo poolstep.Repo
	xactExpRepo  xactexp.Repo
	procExecRepo procexec.Repo
	operator     db.Operator
	log          *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	poolImplRepo Repo,
	poolImplExch Exch,
	poolExecRepo poolexec.Repo,
	implSemRepo implsem.Repo,
	commSemRepo commsem.Repo,
	poolVarRepo poolvar.Repo,
	poolCommRepo poolcomm.Repo,
	poolStepRepo poolstep.Repo,
	xactExpRepo xactexp.Repo,
	procExecRepo procexec.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{
		poolImplRepo, poolImplExch, poolExecRepo, implSemRepo, commSemRepo, poolVarRepo, poolCommRepo,
		poolStepRepo, xactExpRepo, procExecRepo,
		operator, log.With(name),
	}
}

func (s *service) Spawn(spec poolstep.StepSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	refAttr := slog.Any("ref", spec.ImplRef)
	s.log.Debug("proc spawning started", refAttr, slog.Any("exp", spec.PoolExp))
	newRef := implsem.NewRef()
	newImpl := implsem.SemRec{ImplRef: newRef, Kind: implsem.ProcKind}
	newExec := procexec.ExecRec{ImplRef: newRef, LiabMode: implvar.LinearMode}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSemRepo.AddRec(ds, newImpl)
		if err != nil {
			return err
		}
		return s.procExecRepo.InsertRec(ds, newExec)
	})
	if transactErr != nil {
		s.log.Error("proc spawning failed", refAttr)
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("proc spawning succeed", refAttr, slog.Any("proc", newRef))
	return newRef, nil
}

func (s *service) Take(spec poolstep.StepSpec) (err error) {
	ctx := context.Background()
	refAttr := slog.Any("ref", spec.ImplRef)
	s.log.Debug("step taking started", refAttr, slog.Any("exp", spec.PoolExp))
	implSnap, retrErr := s.retrieveSnap(spec.ImplRef)
	if retrErr != nil {
		s.log.Error("step taking failed", refAttr)
		return retrErr
	}
	implMod, implEff, commMod, takeErr := s.take(implSnap, spec.PoolExp)
	if takeErr != nil {
		s.log.Error("step taking failed", refAttr)
		return takeErr
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.poolStepRepo.AddRecs(ds, commMod.Steps)
		if err != nil {
			return err
		}
		err = s.poolCommRepo.ModifyRec(ds, commMod)
		if err != nil {
			return err
		}
		err = s.commSemRepo.TouchRec(ds, commMod.CommRef)
		if err != nil {
			return err
		}
		err = s.poolVarRepo.AddRecs(ds, implMod.Vars)
		if err != nil {
			return err
		}
		return s.implSemRepo.TouchRec(ds, implSnap.ImplRef)
	})
	if transactErr != nil {
		s.log.Error("step taking failed", refAttr)
		return transactErr
	}
	for _, step := range implEff.Steps {
		sendErr := s.poolImplExch.SendSpec(step)
		if sendErr != nil {
			s.log.Error("step taking failed", refAttr)
			return sendErr
		}
	}
	return nil
}

func (s *service) take(
	implSnap ImplSnap,
	exp poolexp.ExpSpec,
) (
	implMod ImplMod,
	implEff ImplEff,
	commMod poolcomm.CommMod,
	err error,
) {
	ctx := context.Background()
	implRefAttr := slog.Any("implRef", implSnap.ImplRef)
	switch poolExp := exp.(type) {
	case poolexp.AcceptSpec:
		commChnl, ok := implSnap.StructVars[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := implSnap.StructExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCtx(poolExp.CommChnlPH)
		}
		nextExpVK := xactExp.(xactexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var commSnap poolcomm.CommSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			commSnap, err = s.poolCommRepo.GetSnapByQry(ds, poolcomm.CommQry{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, getErr
		}
		commMod.CommRef = commSnap.CommRef
		commRefAttr := slog.Any("commRef", commSnap.CommRef)
		subscription := commSnap.NextStep()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступодателя
			implMod.Vars = append(implMod.Vars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем сообщение доступодателя
			commMod.Steps = append(commMod.Steps, poolstep.PubRec{
				CommRef: commSnap.CommRef,
				ImplRef: implSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: poolexp.AcceptRec{
					ContChnlID: newChnlID,
					ContExp:    poolExp.ContExp,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commRefAttr)
			return implMod, implEff, commMod, nil
		}
		acquisition, ok := subscription.(poolstep.SubRec)
		if !ok {
			panic(poolstep.ErrRecTypeUnexpected(subscription))
		}
		switch expRec := acquisition.ContExp.(type) {
		case poolexp.AcquireRec:
			// сдвигаем офсет коммуникации
			commMod.CommON = option.Some(acquisition.CommRef.CommRN)
			// вяжем продолжение доступодателя
			implMod.Vars = append(implMod.Vars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение доступодателя
				implEff.Steps = append(implEff.Steps, poolstep.StepSpec{
					ImplRef: acquisition.ImplRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(poolexp.ErrRecTypeUnexpected(acquisition.ContExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commRefAttr)
		return implMod, implEff, commMod, nil
	case poolexp.AcquireSpec:
		commChnl, ok := implSnap.StructVars[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := implSnap.StructExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCtx(poolExp.CommChnlPH)
		}
		nextExpVK := xactExp.(xactexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var commSnap poolcomm.CommSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			commSnap, err = s.poolCommRepo.GetSnapByQry(ds, poolcomm.CommQry{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, getErr
		}
		commMod.CommRef = commSnap.CommRef
		commAttr := slog.Any("commRef", commSnap.CommRef)
		publication := commSnap.NextStep()
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступополучателя
			implMod.Vars = append(implMod.Vars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем подписку доступополучателя
			commMod.Steps = append(commMod.Steps, poolstep.SubRec{
				CommRef: commSnap.CommRef,
				ImplRef: implSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: poolexp.AcquireRec{
					ContChnlID: newChnlID,
					ContExp:    poolExp.ContExp,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commAttr)
			return implMod, implEff, commMod, nil
		}
		acception, ok := publication.(poolstep.PubRec)
		if !ok {
			panic(poolstep.ErrRecTypeUnexpected(publication))
		}
		switch expRec := acception.ValExp.(type) {
		case poolexp.AcceptRec:
			// сдвигаем офсет коммуникации
			commMod.CommON = option.Some(acception.CommRef.CommRN)
			// вяжем продолжение доступополучателя
			implMod.Vars = append(implMod.Vars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// запускаем продолжение доступополучателя
			implEff.Steps = append(implEff.Steps, poolstep.StepSpec{
				ImplRef: acception.ImplRef,
				PoolExp: expRec.ContExp,
			})
		default:
			panic(poolexp.ErrRecTypeUnexpected(acception.ValExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commAttr)
		return implMod, implEff, commMod, nil
	case poolexp.HireSpec:
		commChnl, ok := implSnap.LinearVars[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := implSnap.LinearExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCtx(poolExp.CommChnlPH)
		}
		nextExpVK := xactExp.(xactexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var commSnap poolcomm.CommSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			commSnap, err = s.poolCommRepo.GetSnapByQry(ds, poolcomm.CommQry{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, getErr
		}
		commMod.CommRef = commSnap.CommRef
		commRefAttr := slog.Any("commRef", commSnap.CommRef)
		subscription := commSnap.NextStep()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение нанимателя
			implMod.Vars = append(implMod.Vars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем сообщение нанимателя
			commMod.Steps = append(commMod.Steps, poolstep.PubRec{
				CommRef: commSnap.CommRef,
				ImplRef: implSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: poolexp.HireRec{
					ContChnlID: newChnlID,
					ProcDescQN: poolExp.ProcDescQN,
					ContExp:    poolExp.ContExp,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commRefAttr)
			return implMod, implEff, commMod, nil
		}
		application, ok := subscription.(poolstep.SubRec)
		if !ok {
			panic(poolstep.ErrRecTypeUnexpected(subscription))
		}
		switch expRec := application.ContExp.(type) {
		case poolexp.ApplyRec:
			// вяжем продолжение нанимателя
			implMod.Vars = append(implMod.Vars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение нанимателя
				implEff.Steps = append(implEff.Steps, poolstep.StepSpec{
					ImplRef: application.ImplRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(poolexp.ErrRecTypeUnexpected(application.ContExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commRefAttr)
		return implMod, implEff, commMod, nil
	case poolexp.ApplySpec:
		commChnl, ok := implSnap.LinearVars[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := implSnap.StructExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCtx(poolExp.CommChnlPH)
		}
		nextExpVK := xactExp.(xactexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var commSnap poolcomm.CommSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			commSnap, err = s.poolCommRepo.GetSnapByQry(ds, poolcomm.CommQry{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, getErr
		}
		commMod.CommRef = commSnap.CommRef
		commAttr := slog.Any("commRef", commSnap.CommRef)
		publication := commSnap.NextStep()
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение соискателя
			implMod.Vars = append(implMod.Vars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем подписку соискателя
			commMod.Steps = append(commMod.Steps, poolstep.SubRec{
				CommRef: commSnap.CommRef,
				ImplRef: implSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: poolexp.ApplyRec{
					ContChnlID: newChnlID,
					ProcDescQN: poolExp.ProcDescQN,
					ContExp:    poolExp.ContExp,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commAttr)
			return implMod, implEff, commMod, nil
		}
		hiring, ok := publication.(poolstep.PubRec)
		if !ok {
			panic(poolstep.ErrRecTypeUnexpected(publication))
		}
		switch expRec := hiring.ValExp.(type) {
		case poolexp.HireRec:
			// вяжем продолжение нанимателя
			implMod.Vars = append(implMod.Vars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// запускаем продолжение нанимателя
			implEff.Steps = append(implEff.Steps, poolstep.StepSpec{
				ImplRef: hiring.ImplRef,
				PoolExp: expRec.ContExp,
			})
		default:
			panic(poolexp.ErrRecTypeUnexpected(hiring.ValExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commAttr)
		return implMod, implEff, commMod, nil
	default:
		panic(poolexp.ErrSpecTypeUnexpected(poolExp))
	}
}

func (s *service) retrieveSnap(ref implsem.SemRef) (_ ImplSnap, err error) {
	ctx := context.Background()
	var implRec ImplRec
	getErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
		implRec, err = s.poolImplRepo.GetRecByRef(ds, ref)
		return err
	})
	if getErr1 != nil {
		return ImplSnap{}, getErr1
	}
	var structExps map[symbol.ADT]xactexp.ExpRec
	var linearExps map[symbol.ADT]xactexp.ExpRec
	getErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
		structExps, err = s.xactExpRepo.GetRecMap(ds, ExtractExpVKs(implRec.StructVars))
		if err != nil {
			return err
		}
		linearExps, err = s.xactExpRepo.GetRecMap(ds, ExtractExpVKs(implRec.LinearVars))
		return err
	})
	if getErr2 != nil {
		return ImplSnap{}, getErr2
	}
	return ImplSnap{
		ImplRef:    implRec.ImplRef,
		StructVars: implvar.ConvertRecsToRecMap(implRec.StructVars),
		StructExps: structExps,
		LinearVars: implvar.ConvertRecsToRecMap(implRec.LinearVars),
		LinearExps: linearExps,
	}, nil
}
