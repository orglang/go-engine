package compexec

import (
	"context"
	"log/slog"

	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/option"
	"orglang/go-engine/adt/poolexec"
	"orglang/go-engine/adt/semcomp"
	"orglang/go-engine/adt/symbol"

	"orglang/go-engine/pool/commexch"
	"orglang/go-engine/pool/commturn"
	"orglang/go-engine/pool/compstep"
	compvar1 "orglang/go-engine/pool/compvar"
	"orglang/go-engine/pool/termexp"
	"orglang/go-engine/pool/typeexp"

	"orglang/go-engine/proc/compexec"
	"orglang/go-engine/proc/termdef"
)

type API interface {
	Take(compstep.StepSpec) error
	Spawn(compstep.StepSpec) (semcomp.CompRef, error)
}

type ExecRec struct {
	CompRef    semcomp.CompRef
	LiabMode   compvar.Mode
	StructVars []compvar.StructRec
	LinearVars []compvar.LinearRec
}

type ExecSnap struct {
	CompRef    semcomp.CompRef
	StructVars map[symbol.ADT]compvar.StructRec
	StructExps map[symbol.ADT]typeexp.ExpRec
	LinearVars map[symbol.ADT]compvar.LinearRec
	LinearExps map[symbol.ADT]typeexp.ExpRec
}

type ExecMod struct {
	Vars []compvar.VarRec
}

type ExecEff struct {
	Steps []compstep.StepSpec
}

type service struct {
	compExecRepo Repo
	poolImplExch Exch
	poolExecRepo poolexec.Repo
	compVarRepo  compvar1.Repo
	commExchRepo commexch.Repo
	commTurnRepo commturn.Repo
	xactExpRepo  typeexp.Repo
	procExecRepo compexec.Repo
	operator     db.Operator
	log          *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	compExecRepo Repo,
	compExecExch Exch,
	poolExecRepo poolexec.Repo,
	poolVarRepo compvar1.Repo,
	poolCommRepo commexch.Repo,
	poolStepRepo commturn.Repo,
	xactExpRepo typeexp.Repo,
	procExecRepo compexec.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{
		compExecRepo, compExecExch, poolExecRepo, poolVarRepo,
		poolCommRepo, poolStepRepo, xactExpRepo, procExecRepo,
		operator, log.With(name),
	}
}

func (s *service) Spawn(spec compstep.StepSpec) (_ semcomp.CompRef, err error) {
	ctx := context.Background()
	refAttr := slog.Any("ref", spec.CompRef)
	s.log.Debug("proc spawning started", refAttr, slog.Any("exp", spec.PoolExp))
	newExec := compexec.ExecRec{CompRef: semcomp.NewRef(), LiabMode: compvar.LinearMode}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		return s.procExecRepo.InsertRec(ds, newExec)
	})
	if transactErr != nil {
		s.log.Error("proc spawning failed", refAttr)
		return semcomp.CompRef{}, transactErr
	}
	s.log.Debug("proc spawning succeed", refAttr, slog.Any("proc", newExec.CompRef))
	return newExec.CompRef, nil
}

func (s *service) Take(spec compstep.StepSpec) (err error) {
	ctx := context.Background()
	refAttr := slog.Any("ref", spec.CompRef)
	s.log.Debug("step taking started", refAttr, slog.Any("exp", spec.PoolExp))
	execSnap, retrErr := s.retrieveSnap(spec.CompRef)
	if retrErr != nil {
		s.log.Error("step taking failed", refAttr)
		return retrErr
	}
	execMod, execEff, exchMod, takeErr := s.take(execSnap, spec.PoolExp)
	if takeErr != nil {
		s.log.Error("step taking failed", refAttr)
		return takeErr
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.commTurnRepo.AddRecs(ds, exchMod.Turns)
		if err != nil {
			return err
		}
		err = s.commExchRepo.ModifyRec(ds, exchMod)
		if err != nil {
			return err
		}
		err = s.compVarRepo.AddRecs(ds, execMod.Vars)
		if err != nil {
			return err
		}
		return s.compExecRepo.ModifyRec(ds, execMod)
	})
	if transactErr != nil {
		s.log.Error("step taking failed", refAttr)
		return transactErr
	}
	for _, step := range execEff.Steps {
		sendErr := s.poolImplExch.SendSpec(step)
		if sendErr != nil {
			s.log.Error("step taking failed", refAttr)
			return sendErr
		}
	}
	return nil
}

func (s *service) take(
	execSnap ExecSnap,
	exp termexp.ExpSpec,
) (
	execMod ExecMod,
	execEff ExecEff,
	exchMod commexch.ExchMod,
	err error,
) {
	ctx := context.Background()
	implRefAttr := slog.Any("implRef", execSnap.CompRef)
	switch poolExp := exp.(type) {
	case termexp.AcceptSpec:
		commChnl, ok := execSnap.StructVars[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.StructExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCtx(poolExp.CommChnlPH)
		}
		nextExpVK := xactExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var commSnap commexch.ExchSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			commSnap, err = s.commExchRepo.GetSnapByQry(ds, commexch.ExchQry{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, getErr
		}
		exchMod.CommRef = commSnap.CommRef
		commRefAttr := slog.Any("commRef", commSnap.CommRef)
		subscription := commSnap.NextTurn()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступодателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.CompRef,
				ExchRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем сообщение доступодателя
			exchMod.Turns = append(exchMod.Turns, commturn.PubRec{
				CommRef: commSnap.CommRef,
				CompRef: execSnap.CompRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: termexp.AcceptRec{
					ContChnlID: newChnlID,
					ContExp:    poolExp.ContExp,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commRefAttr)
			return execMod, execEff, exchMod, nil
		}
		acquisition, ok := subscription.(commturn.SubRec)
		if !ok {
			panic(commturn.ErrRecTypeUnexpected(subscription))
		}
		if poolExp.ContExp != nil {
			// шедулим продолжение доступополучателя
			execEff.Steps = append(execEff.Steps, compstep.StepSpec{
				CompRef: execSnap.CompRef,
				PoolExp: poolExp.ContExp,
			})
		}
		switch expRec := acquisition.ContExp.(type) {
		case termexp.AcquireRec:
			// сдвигаем офсет коммуникации
			exchMod.OffsetNr = option.Some(acquisition.CommRef.CommRN)
			// вяжем продолжение доступодателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.CompRef,
				ExchRef: commChnl.CommRef,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение доступодателя
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: acquisition.CompRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(termexp.ErrRecTypeUnexpected(acquisition.ContExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commRefAttr)
		return execMod, execEff, exchMod, nil
	case termexp.AcquireSpec:
		commChnl, ok := execSnap.StructVars[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.StructExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCtx(poolExp.CommChnlPH)
		}
		nextExpVK := xactExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var commSnap commexch.ExchSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			commSnap, err = s.commExchRepo.GetSnapByQry(ds, commexch.ExchQry{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, getErr
		}
		exchMod.CommRef = commSnap.CommRef
		commAttr := slog.Any("commRef", commSnap.CommRef)
		publication := commSnap.NextTurn()
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступополучателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.CompRef,
				ExchRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем подписку доступополучателя
			exchMod.Turns = append(exchMod.Turns, commturn.SubRec{
				CommRef: commSnap.CommRef,
				CompRef: execSnap.CompRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: termexp.AcquireRec{
					ContChnlID: newChnlID,
					ContExp:    poolExp.ContExp,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commAttr)
			return execMod, execEff, exchMod, nil
		}
		acception, ok := publication.(commturn.PubRec)
		if !ok {
			panic(commturn.ErrRecTypeUnexpected(publication))
		}
		if poolExp.ContExp != nil {
			// шедулим продолжение доступодателя
			execEff.Steps = append(execEff.Steps, compstep.StepSpec{
				CompRef: execSnap.CompRef,
				PoolExp: poolExp.ContExp,
			})
		}
		switch expRec := acception.ValExp.(type) {
		case termexp.AcceptRec:
			// сдвигаем офсет коммуникации
			exchMod.OffsetNr = option.Some(acception.CommRef.CommRN)
			// вяжем продолжение доступополучателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.CompRef,
				ExchRef: commChnl.CommRef,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение доступополучателя
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: acception.CompRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(termexp.ErrRecTypeUnexpected(acception.ValExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commAttr)
		return execMod, execEff, exchMod, nil
	case termexp.ApplySpec:
		commChnl, ok := execSnap.LinearVars[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.LinearExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCtx(poolExp.CommChnlPH)
		}
		nextExpVK := xactExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var commSnap commexch.ExchSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			commSnap, err = s.commExchRepo.GetSnapByQry(ds, commexch.ExchQry{
				CommRef: commChnl.ExchRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, getErr
		}
		exchMod.CommRef = commSnap.CommRef
		commRefAttr := slog.Any("commRef", commSnap.CommRef)
		subscription := commSnap.NextTurn()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение соискателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.TermRef,
				ExchRef: commChnl.ExchRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем сообщение соискателя
			exchMod.Turns = append(exchMod.Turns, commturn.PubRec{
				CommRef: commSnap.CommRef,
				CompRef: execSnap.CompRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: termexp.ApplyRec{
					ContChnlID: newChnlID,
					ProcDescQN: poolExp.ProcDescQN,
					ContExp:    poolExp.ContExp,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commRefAttr)
			return execMod, execEff, exchMod, nil
		}
		hiring, ok := subscription.(commturn.SubRec)
		if !ok {
			panic(commturn.ErrRecTypeUnexpected(subscription))
		}
		if poolExp.ContExp != nil {
			// шедулим продолжение нанимателя
			execEff.Steps = append(execEff.Steps, compstep.StepSpec{
				CompRef: execSnap.CompRef,
				PoolExp: poolExp.ContExp,
			})
		}
		switch expRec := hiring.ContExp.(type) {
		case termexp.HireRec:
			// вяжем продолжение соискателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.TermRef,
				ExchRef: commChnl.ExchRef,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение соискателя
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: hiring.CompRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(termexp.ErrRecTypeUnexpected(hiring.ContExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commRefAttr)
		return execMod, execEff, exchMod, nil
	case termexp.HireSpec:
		commChnl, ok := execSnap.LinearVars[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.LinearExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCtx(poolExp.CommChnlPH)
		}
		nextExpVK := xactExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var commSnap commexch.ExchSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			commSnap, err = s.commExchRepo.GetSnapByQry(ds, commexch.ExchQry{
				CommRef: commChnl.ExchRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, getErr
		}
		exchMod.CommRef = commSnap.CommRef
		commAttr := slog.Any("commRef", commSnap.CommRef)
		publication := commSnap.NextTurn()
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение нанимателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.TermRef,
				ExchRef: commChnl.ExchRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем подписку нанимателя
			exchMod.Turns = append(exchMod.Turns, commturn.SubRec{
				CommRef: commSnap.CommRef,
				CompRef: execSnap.CompRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: termexp.HireRec{
					ContChnlID: newChnlID,
					ProcDescQN: poolExp.ProcDescQN,
					ContExp:    poolExp.ContExp,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commAttr)
			return execMod, execEff, exchMod, nil
		}
		application, ok := publication.(commturn.PubRec)
		if !ok {
			panic(commturn.ErrRecTypeUnexpected(publication))
		}
		if poolExp.ContExp != nil {
			// шедулим продолжение соискателя
			execEff.Steps = append(execEff.Steps, compstep.StepSpec{
				CompRef: execSnap.CompRef,
				PoolExp: poolExp.ContExp,
			})
		}
		switch expRec := application.ValExp.(type) {
		case termexp.ApplyRec:
			// вяжем продолжение нанимателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.TermRef,
				ExchRef: commChnl.ExchRef,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			if expRec.ContExp != nil {
				// запускаем продолжение нанимателя
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: application.CompRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(termexp.ErrRecTypeUnexpected(application.ValExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commAttr)
		return execMod, execEff, exchMod, nil
	case termexp.ReleaseSpec:
		commChnl, ok := execSnap.LinearVars[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.LinearExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCtx(poolExp.CommChnlPH)
		}
		nextExpVK := xactExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var commSnap commexch.ExchSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			commSnap, err = s.commExchRepo.GetSnapByQry(ds, commexch.ExchQry{
				CommRef: commChnl.ExchRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, getErr
		}
		exchMod.CommRef = commSnap.CommRef
		commRefAttr := slog.Any("commRef", commSnap.CommRef)
		subscription := commSnap.NextTurn()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступовозвращателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.TermRef,
				ExchRef: commChnl.ExchRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем сообщение доступовозвращателя
			exchMod.Turns = append(exchMod.Turns, commturn.PubRec{
				CommRef: commSnap.CommRef,
				CompRef: execSnap.CompRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: termexp.ReleaseRec{
					ContChnlID: newChnlID,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commRefAttr)
			return execMod, execEff, exchMod, nil
		}
		detaching, ok := subscription.(commturn.SubRec)
		if !ok {
			panic(commturn.ErrRecTypeUnexpected(subscription))
		}
		switch expRec := detaching.ContExp.(type) {
		case termexp.DetachRec:
			// вяжем продолжение доступовозвращателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.TermRef,
				ExchRef: commChnl.ExchRef,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
		default:
			panic(termexp.ErrRecTypeUnexpected(detaching.ContExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commRefAttr)
		return execMod, execEff, exchMod, nil
	case termexp.DetachSpec:
		commChnl, ok := execSnap.LinearVars[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.LinearExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef.ErrMissingInCtx(poolExp.CommChnlPH)
		}
		nextExpVK := xactExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var commSnap commexch.ExchSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			commSnap, err = s.commExchRepo.GetSnapByQry(ds, commexch.ExchQry{
				CommRef: commChnl.ExchRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, getErr
		}
		exchMod.CommRef = commSnap.CommRef
		commAttr := slog.Any("commRef", commSnap.CommRef)
		publication := commSnap.NextTurn()
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение доступопринимателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.TermRef,
				ExchRef: commChnl.ExchRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем подписку доступопринимателя
			exchMod.Turns = append(exchMod.Turns, commturn.SubRec{
				CommRef: commSnap.CommRef,
				CompRef: execSnap.CompRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: termexp.DetachRec{
					ContChnlID: newChnlID,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commAttr)
			return execMod, execEff, exchMod, nil
		}
		releasing, ok := publication.(commturn.PubRec)
		if !ok {
			panic(commturn.ErrRecTypeUnexpected(publication))
		}
		switch expRec := releasing.ValExp.(type) {
		case termexp.ReleaseRec:
			// вяжем продолжение доступопринимателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				TermRef: commChnl.TermRef,
				ExchRef: commChnl.ExchRef,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
		default:
			panic(termexp.ErrRecTypeUnexpected(releasing.ValExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commAttr)
		return execMod, execEff, exchMod, nil
	default:
		panic(termexp.ErrSpecTypeUnexpected(exp))
	}
}

func (s *service) retrieveSnap(ref semcomp.CompRef) (_ ExecSnap, err error) {
	ctx := context.Background()
	var execRec ExecRec
	getErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
		execRec, err = s.compExecRepo.GetRecByRef(ds, ref)
		return err
	})
	if getErr1 != nil {
		return ExecSnap{}, getErr1
	}
	var structExps map[symbol.ADT]typeexp.ExpRec
	var linearExps map[symbol.ADT]typeexp.ExpRec
	getErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
		structExps, err = s.xactExpRepo.GetRecMap(ds, ExtractExpVKs(execRec.StructVars))
		if err != nil {
			return err
		}
		linearExps, err = s.xactExpRepo.GetRecMap(ds, ExtractExpVKs(execRec.LinearVars))
		return err
	})
	if getErr2 != nil {
		return ExecSnap{}, getErr2
	}
	return ExecSnap{
		CompRef:    execRec.CompRef,
		StructVars: compvar.ConvertRecsToRecMap(execRec.StructVars),
		StructExps: structExps,
		LinearVars: compvar.ConvertRecsToRecMap(execRec.LinearVars),
		LinearExps: linearExps,
	}, nil
}
