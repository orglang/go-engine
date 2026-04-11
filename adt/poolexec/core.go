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
	"orglang/go-engine/adt/poolcomm"
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
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Run(ExecSpec) (implsem.SemRef, error) // aka Create
	Take(poolstep.StepSpec) error
	Spawn(poolstep.StepSpec) (implsem.SemRef, error)
	RetrieveSnap(implsem.SemRef) (ExecCfgSnap, error)
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

type ImplMod struct {
	ImplRef implsem.SemRef
	Vars    []implvar.VarRec
}

type ImplEff struct {
	PoolSteps []poolstep.StepSpec
}

type ExecCfgSnap struct {
	ImplRef    implsem.SemRef
	StructVars map[symbol.ADT]implvar.StructRec
	LinearVars map[symbol.ADT]implvar.LinearRec
}

type ExecCtxSnap struct {
	StructExps map[symbol.ADT]xactexp.ExpRec
	LinearExps map[symbol.ADT]xactexp.ExpRec
}

type ExecLiabSnap struct {
	ImplRef implsem.SemRef
	LiabVar implvar.VarRec
}

type service struct {
	poolExecRepo Repo
	implSemRepo  implsem.Repo
	commSemRepo  commsem.Repo
	poolDecRepo  pooldec.Repo
	poolVarRepo  poolvar.Repo
	poolCommRepo poolcomm.Repo
	poolStepRepo poolstep.Repo
	xactDefRepo  xactexp.Repo
	xactExpRepo  xactexp.Repo
	procDecRepo  procdec.Repo
	procExecRepo procexec.Repo
	poolStepExch poolstep.Exch
	operator     db.Operator
	log          *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	poolExecRepo Repo,
	implSemRepo implsem.Repo,
	commSemRepo commsem.Repo,
	poolDecRepo pooldec.Repo,
	poolVarRepo poolvar.Repo,
	poolConnRepo poolcomm.Repo,
	poolStepRepo poolstep.Repo,
	xactDefRepo xactexp.Repo,
	xactExpRepo xactexp.Repo,
	procDecRepo procdec.Repo,
	procExecRepo procexec.Repo,
	poolStepPub poolstep.Exch,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{
		poolExecRepo, implSemRepo, commSemRepo, poolDecRepo, poolVarRepo, poolConnRepo,
		poolStepRepo, xactDefRepo, xactExpRepo, procDecRepo, procExecRepo,
		poolStepPub, operator, log.With(name),
	}
}

func (s *service) Run(spec ExecSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	specAttr := slog.Any("spec", spec)
	s.log.Debug("creation started", specAttr)
	var poolDec pooldec.DecRec
	getErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
		poolDec, err = s.poolDecRepo.GetRecByQN(ds, spec.DescQN)
		return err
	})
	if getErr1 != nil {
		s.log.Error("creation failed", specAttr)
		return implsem.SemRef{}, getErr1
	}
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars))
	for _, assetVar := range spec.AssetVars {
		if assetVar.ImplQN == spec.LiabVar.ImplQN {
			continue
		}
		assetQNs = append(assetQNs, assetVar.ImplQN)
	}
	var assetExecs map[uniqsym.ADT]ExecLiabSnap
	getErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
		assetExecs, err = s.poolExecRepo.GetSnapsByQNs(ds, assetQNs)
		return err
	})
	if getErr2 != nil {
		s.log.Error("creation failed", specAttr)
		return implsem.SemRef{}, getErr2
	}
	newImplSem := implsem.SemRec{ImplRef: implsem.NewRef(), ImplQN: spec.LiabVar.ImplQN, Kind: implsem.Pool}
	newCommSem := commsem.SemRec{CommRef: commsem.NewRef(), Kind: commsem.Pool}
	newConn := poolcomm.ConnRec{CommRef: newCommSem.CommRef, CommON: newCommSem.CommRef.CommRN}
	newExec := ExecRec{ImplRef: newImplSem.ImplRef, LiabMode: implvar.StructMode}
	newLiabVar := implvar.StructRec{
		ImplRef: newImplSem.ImplRef,
		CommRef: newCommSem.CommRef,
		ChnlID:  identity.New(),
		ChnlPH:  spec.LiabVar.ChnlPH,
		ChnlBS:  implvar.LiabSide,
		ExpVK:   poolDec.LiabVar.ExpVK,
	}
	newAssetVars := make([]implvar.VarRec, 0, len(spec.AssetVars)+1)
	for _, assetVar := range spec.AssetVars {
		var commRef commsem.SemRef
		var chnlID identity.ADT
		var expVK valkey.ADT
		assetExec, ok := assetExecs[assetVar.ImplQN]
		if !ok && assetVar.ImplQN == spec.LiabVar.ImplQN {
			commRef = newCommSem.CommRef
			chnlID = newLiabVar.ChnlID
			expVK = newLiabVar.ExpVK
		} else {
			commRef = assetExec.LiabVar.GetCommRef()
			chnlID = assetExec.LiabVar.GetChnlID()
			expVK = assetExec.LiabVar.GetExpVK()
		}
		newAssetVars = append(newAssetVars, implvar.StructRec{
			ImplRef: newImplSem.ImplRef,
			CommRef: commRef,
			ChnlID:  chnlID,
			ChnlPH:  assetVar.ChnlPH,
			ChnlBS:  implvar.AssetSide,
			ExpVK:   expVK,
		})
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSemRepo.AddRec(ds, newImplSem)
		if err != nil {
			return err
		}
		err = s.poolExecRepo.AddRec(ds, newExec)
		if err != nil {
			return err
		}
		err = s.poolVarRepo.AddRecs(ds, append(newAssetVars, newLiabVar))
		if err != nil {
			return err
		}
		err = s.commSemRepo.AddRec(ds, newCommSem)
		if err != nil {
			return err
		}
		return s.poolCommRepo.AddRec(ds, newConn)
	})
	if transactErr != nil {
		s.log.Error("creation failed", specAttr)
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", slog.Any("ref", newImplSem.ImplRef))
	return newImplSem.ImplRef, nil
}

func (s *service) RetrieveSnap(ref implsem.SemRef) (snap ExecCfgSnap, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		snap, err = s.poolExecRepo.GetSnap(ds, ref)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("ref", ref))
		return ExecCfgSnap{}, err
	}
	return snap, nil
}

func (s *service) Spawn(spec poolstep.StepSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	refAttr := slog.Any("ref", spec.ImplRef)
	s.log.Debug("proc spawning started", refAttr, slog.Any("exp", spec.CommExp))
	newRef := implsem.NewRef()
	newImpl := implsem.SemRec{ImplRef: newRef, Kind: implsem.Proc}
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
	s.log.Debug("step taking started", refAttr, slog.Any("exp", spec.CommExp))
	var execCfg ExecCfgSnap
	selectErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
		execCfg, err = s.poolExecRepo.GetSnap(ds, spec.ImplRef)
		return err
	})
	if selectErr1 != nil {
		s.log.Error("step taking failed", refAttr)
		return selectErr1
	}
	var execCtx ExecCtxSnap
	selectErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
		execCtx.StructExps, err = s.xactExpRepo.GetRecMap(ds, ExtractExpVKs(execCfg.StructVars))
		if err != nil {
			return err
		}
		execCtx.LinearExps, err = s.xactExpRepo.GetRecMap(ds, ExtractExpVKs(execCfg.LinearVars))
		return err
	})
	if selectErr2 != nil {
		s.log.Error("step taking failed", refAttr)
		return selectErr2
	}
	execMod, implEff, connMod, takeErr := s.take(execCfg, execCtx, spec.CommExp)
	if takeErr != nil {
		s.log.Error("step taking failed", refAttr)
		return takeErr
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.poolStepRepo.AddRecs(ds, connMod.Steps)
		if err != nil {
			return err
		}
		err = s.poolCommRepo.ModifyRec(ds, connMod)
		if err != nil {
			return err
		}
		err = s.commSemRepo.TouchRec(ds, connMod.CommRef)
		if err != nil {
			return err
		}
		err = s.poolVarRepo.AddRecs(ds, execMod.Vars)
		if err != nil {
			return err
		}
		return s.implSemRepo.TouchRec(ds, execMod.ImplRef)
	})
	if transactErr != nil {
		s.log.Error("step taking failed", refAttr)
		return transactErr
	}
	for _, step := range implEff.PoolSteps {
		sendErr := s.poolStepExch.SendSpec(step)
		if sendErr != nil {
			s.log.Error("step taking failed", refAttr)
			return sendErr
		}
	}
	s.log.Debug("step taking succeed", refAttr)
	return nil
}

func (s *service) take(
	execCfg ExecCfgSnap,
	execCtx ExecCtxSnap,
	exp poolexp.ExpSpec,
) (
	implMod ImplMod,
	implEff ImplEff,
	commMod poolcomm.CommMod,
	err error,
) {
	ctx := context.Background()
	implMod.ImplRef = execCfg.ImplRef
	implRefAttr := slog.Any("implRef", execCfg.ImplRef)
	switch commExp := exp.(type) {
	case poolexp.AcceptSpec:
		commChnl, ok := execCfg.StructVars[commExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCfg(commExp.CommChnlPH)
		}
		// вычисляем следующее состояние
		xactExp, ok := execCtx.StructExps[commExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCtx(commExp.CommChnlPH)
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
				ChnlID:  commChnl.ChnlID,
				ValExp: poolexp.AcceptRec{
					ContChnlID: newChnlID,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commRefAttr)
			return implMod, implEff, commMod, nil
		}
		acquisition, ok := subscription.(poolstep.SubRec)
		if !ok {
			panic(poolstep.ErrRecTypeUnexpected(subscription))
		}
		switch contExp := acquisition.ContExp.(type) {
		case poolexp.AcquireRec:
			// сдвигаем офсет коммуникации
			commMod.CommON = option.Some(acquisition.CommRef.CommRN)
			// вяжем продолжение доступодателя
			implMod.Vars = append(implMod.Vars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// запускаем продолжение доступодателя
			implEff.PoolSteps = append(implEff.PoolSteps, poolstep.StepSpec{
				ImplRef: execCfg.ImplRef,
				CommExp: contExp.ContExp,
			})
		default:
			panic(poolexp.ErrRecTypeUnexpected(acquisition.ContExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commRefAttr)
		return implMod, implEff, commMod, nil
	case poolexp.AcquireSpec:
		commChnl, ok := execCfg.StructVars[commExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCfg(commExp.CommChnlPH)
		}
		xactExp, ok := execCtx.StructExps[commExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCtx(commExp.CommChnlPH)
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
				ChnlID:  commChnl.ChnlID,
				ContExp: poolexp.AcquireRec{
					ContChnlID: newChnlID,
				},
			})
			s.log.Debug("taking half done", implRefAttr, commAttr)
			return implMod, implEff, commMod, nil
		}
		acception, ok := publication.(poolstep.PubRec)
		if !ok {
			panic(poolstep.ErrRecTypeUnexpected(publication))
		}
		switch valExp := acception.ValExp.(type) {
		case poolexp.AcceptRec:
			// сдвигаем офсет коммуникации
			commMod.CommON = option.Some(acception.CommRef.CommRN)
			// вяжем продолжение доступополучателя
			implMod.Vars = append(implMod.Vars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// запускаем продолжение доступополучателя
			implEff.PoolSteps = append(implEff.PoolSteps, poolstep.StepSpec{
				ImplRef: execCfg.ImplRef,
				CommExp: valExp.ContExp,
			})
		default:
			panic(poolexp.ErrRecTypeUnexpected(acception.ValExp))
		}
		s.log.Debug("step taking succeed", implRefAttr, commAttr)
		return implMod, implEff, commMod, nil
	case poolexp.HireSpec:
		commChnl, ok := execCfg.LinearVars[commExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCfg(commExp.CommChnlPH)
		}
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
		}

		// ...

		s.log.Debug("step taking succeed", implRefAttr, commRefAttr)
		return implMod, implEff, commMod, nil
	case poolexp.ApplySpec:
		commChnl, ok := execCfg.LinearVars[commExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implRefAttr)
			return implMod, implEff, commMod, procdef.ErrMissingInCfg(commExp.CommChnlPH)
		}
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

		}

		// ...

		s.log.Debug("step taking succeed", implRefAttr, commAttr)
		return implMod, implEff, commMod, nil
	default:
		panic(poolexp.ErrSpecTypeUnexpected(commExp))
	}
}
