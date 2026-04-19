package compexec

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/option"
	"orglang/go-engine/adt/seqnum"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"

	"orglang/go-engine/pool/commexch"
	"orglang/go-engine/pool/commturn"
	"orglang/go-engine/pool/compstep"
	"orglang/go-engine/pool/compvar"
	"orglang/go-engine/pool/termdef"
	"orglang/go-engine/pool/termexp"
	"orglang/go-engine/pool/typeexp"

	"orglang/go-engine/proc/compexec"
	termdef1 "orglang/go-engine/proc/termdef"
)

type API interface {
	Run(ExecSpec) (compsem.SemRef, error) // aka Create
	Take(compstep.StepSpec) error
	Spawn(compstep.StepSpec) (compsem.SemRef, error)
}

type ExecSpec struct {
	// ссылка на декларацию вновь создаваемого пула
	TermQN uniqsym.ADT
	// внутренняя и внешняя ссылки на вновь создаваемый пул
	LiabVar compvar.VarSpec
	// внутренние и внешние ссылки на ранее созданные пулы
	AssetVars []compvar.VarSpec
}

type ExecRec struct {
	CompRef  compsem.SemRef
	LiabMode compvar.Mode
}

type ExecMod struct {
	CompRef compsem.SemRef
	Vars    []compvar.VarRec
}

func (mod ExecMod) isEmpty() bool { return len(mod.Vars) == 0 }

type ExecEff struct {
	Steps []compstep.StepSpec
}

type ExecSnap1 struct {
	CompRef compsem.SemRef
	LiabVar compvar.VarRec
}

type ExecSnap2 struct {
	CompRef    compsem.SemRef
	LiabMode   compvar.Mode
	StructVars []compvar.StructRec
	LinearVars []compvar.LinearRec
}

type ExecSnap3 struct {
	CompRef    compsem.SemRef
	StructVars map[symbol.ADT]compvar.StructRec
	StructExps map[symbol.ADT]typeexp.ExpRec
	LinearVars map[symbol.ADT]compvar.LinearRec
	LinearExps map[symbol.ADT]typeexp.ExpRec
}

type service struct {
	compExecRepo   Repo
	compExecBroker Broker
	compVarRepo    compvar.Repo
	commExchRepo   commexch.Repo
	commTurnRepo   commturn.Repo
	typeExpRepo    typeexp.Repo
	procExecRepo   compexec.Repo
	termDefRepo    termdef.Repo
	implSemRepo    implsem.Repo
	compSemRepo    compsem.Repo
	operator       db.Operator
	log            *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	compExecRepo Repo,
	compExecExch Broker,
	compVarRepo compvar.Repo,
	commExchRepo commexch.Repo,
	commTurnRepo commturn.Repo,
	typeExpRepo typeexp.Repo,
	procExecRepo compexec.Repo,
	termDefRepo termdef.Repo,
	implSemRepo implsem.Repo,
	compSemRepo compsem.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{
		compExecRepo, compExecExch, compVarRepo,
		commExchRepo, commTurnRepo, typeExpRepo, procExecRepo, termDefRepo,
		implSemRepo, compSemRepo,
		operator, log.With(name),
	}
}

func (s *service) Run(spec ExecSpec) (_ compsem.SemRef, err error) {
	ctx := context.Background()
	specAttr := slog.Any("spec", spec)
	s.log.Debug("creation started", specAttr)
	var termDec termdef.DefRec
	getErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
		termDec, err = s.termDefRepo.GetRecByQN(ds, spec.TermQN)
		return err
	})
	if getErr1 != nil {
		s.log.Error("creation failed", specAttr)
		return compsem.SemRef{}, getErr1
	}
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars))
	for _, assetVar := range spec.AssetVars {
		if assetVar.TermQN == spec.LiabVar.TermQN {
			continue
		}
		assetQNs = append(assetQNs, assetVar.TermQN)
	}
	var assetExecs map[uniqsym.ADT]ExecSnap1
	getErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
		assetExecs, err = s.compExecRepo.GetSnapMapByQNs(ds, assetQNs)
		return err
	})
	if getErr2 != nil {
		s.log.Error("creation failed", specAttr)
		return compsem.SemRef{}, getErr2
	}
	newExec := ExecRec{CompRef: compsem.New(), LiabMode: compvar.StructMode}
	newExch := commexch.ExchRec{CommRef: commsem.New(), OffsetNr: seqnum.Zero}
	newImpl := implsem.SemRec{ImplQN: spec.TermQN, ImplID: newExec.CompRef.CompID}
	newLiabVar := compvar.StructRec{
		CompRef: newExec.CompRef,
		CommRef: newExch.CommRef,
		ChnlID:  identity.New(),
		ChnlPH:  spec.LiabVar.ChnlPH,
		ChnlBS:  compvar.LiabSide,
		ExpVK:   termDec.LiabVar.ExpVK,
	}
	newAssetVars := make([]compvar.VarRec, 0, len(spec.AssetVars)+1)
	for _, assetVar := range spec.AssetVars {
		var commRef commsem.SemRef
		var chnlID identity.ADT
		var expVK valkey.ADT
		assetExec, ok := assetExecs[assetVar.TermQN]
		if !ok && assetVar.TermQN == spec.LiabVar.TermQN {
			commRef = newExch.CommRef
			chnlID = newLiabVar.ChnlID
			expVK = newLiabVar.ExpVK
		} else {
			commRef = assetExec.LiabVar.GetCommRef()
			chnlID = assetExec.LiabVar.GetChnlID()
			expVK = assetExec.LiabVar.GetExpVK()
		}
		newAssetVars = append(newAssetVars, compvar.StructRec{
			CompRef: newExec.CompRef,
			CommRef: commRef,
			ChnlID:  chnlID,
			ChnlPH:  assetVar.ChnlPH,
			ChnlBS:  compvar.AssetSide,
			ExpVK:   expVK,
		})
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSemRepo.AddRec(ds, newImpl)
		if err != nil {
			return err
		}
		err = s.compExecRepo.AddRec(ds, newExec)
		if err != nil {
			return err
		}
		err = s.compVarRepo.AddRecs(ds, append(newAssetVars, newLiabVar))
		if err != nil {
			return err
		}
		return s.commExchRepo.AddRec(ds, newExch)
	})
	if transactErr != nil {
		s.log.Error("creation failed", specAttr)
		return compsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", slog.Any("ref", newExec.CompRef))
	return newExec.CompRef, nil
}

func (s *service) Spawn(spec compstep.StepSpec) (_ compsem.SemRef, err error) {
	ctx := context.Background()
	refAttr := slog.Any("ref", spec.CompRef)
	s.log.Debug("proc spawning started", refAttr, slog.Any("exp", spec.PoolExp))
	newExec := compexec.ExecRec{CompRef: compsem.New(), LiabMode: compvar.LinearMode}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		return s.procExecRepo.AddRec(ds, newExec)
	})
	if transactErr != nil {
		s.log.Error("proc spawning failed", refAttr)
		return compsem.SemRef{}, transactErr
	}
	s.log.Debug("proc spawning succeed", refAttr, slog.Any("proc", newExec.CompRef))
	return newExec.CompRef, nil
}

func (s *service) Take(spec compstep.StepSpec) (err error) {
	ctx := context.Background()
	refAttr := slog.Any("ref", spec.CompRef)
	s.log.Debug("step taking started", refAttr, slog.Any("exp", spec.PoolExp))
	execSnap, retErr := s.retrieveSnap(spec.CompRef)
	if retErr != nil {
		s.log.Error("step taking failed", refAttr)
		return retErr
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
		return s.compSemRepo.TouchRef(ds, execSnap.CompRef)
	})
	if transactErr != nil {
		s.log.Error("step taking failed", refAttr)
		return transactErr
	}
	for _, step := range execEff.Steps {
		sendErr := s.compExecBroker.SendSpec(step)
		if sendErr != nil {
			s.log.Error("step taking failed", refAttr)
			return sendErr
		}
	}
	return nil
}

func (s *service) take(
	execSnap ExecSnap3,
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
			return execMod, execEff, exchMod, termdef1.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.StructExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef1.ErrMissingInCtx(poolExp.CommChnlPH)
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
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
			return execMod, execEff, exchMod, termdef1.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.StructExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef1.ErrMissingInCtx(poolExp.CommChnlPH)
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
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
			return execMod, execEff, exchMod, termdef1.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.LinearExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef1.ErrMissingInCtx(poolExp.CommChnlPH)
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
			// вяжем продолжение соискателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
					ProcTermQN: poolExp.ProcTermQN,
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
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
			return execMod, execEff, exchMod, termdef1.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.LinearExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef1.ErrMissingInCtx(poolExp.CommChnlPH)
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
			// вяжем продолжение нанимателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
					ProcTermQN: poolExp.ProcTermQN,
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
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
			return execMod, execEff, exchMod, termdef1.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.LinearExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef1.ErrMissingInCtx(poolExp.CommChnlPH)
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
			// вяжем продолжение доступовозвращателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
			return execMod, execEff, exchMod, termdef1.ErrMissingInCfg(poolExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execSnap.LinearExps[poolExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return execMod, execEff, exchMod, termdef1.ErrMissingInCtx(poolExp.CommChnlPH)
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
			// вяжем продолжение доступопринимателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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
				CompRef: commChnl.CompRef,
				CommRef: commChnl.CommRef,
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

func (s *service) retrieveSnap(ref compsem.SemRef) (_ ExecSnap3, err error) {
	ctx := context.Background()
	var execSnap ExecSnap2
	getErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
		execSnap, err = s.compExecRepo.GetSnapByRef(ds, ref)
		return err
	})
	if getErr1 != nil {
		return ExecSnap3{}, getErr1
	}
	var structExps map[symbol.ADT]typeexp.ExpRec
	var linearExps map[symbol.ADT]typeexp.ExpRec
	getErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
		structExps, err = s.typeExpRepo.GetRecMap(ds, ExtractExpVKs(execSnap.StructVars))
		if err != nil {
			return err
		}
		linearExps, err = s.typeExpRepo.GetRecMap(ds, ExtractExpVKs(execSnap.LinearVars))
		return err
	})
	if getErr2 != nil {
		return ExecSnap3{}, getErr2
	}
	return ExecSnap3{
		CompRef:    execSnap.CompRef,
		StructVars: compvar.ConvertRecsToRecMap(execSnap.StructVars),
		StructExps: structExps,
		LinearVars: compvar.ConvertRecsToRecMap(execSnap.LinearVars),
		LinearExps: linearExps,
	}, nil
}
