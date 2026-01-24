package procexec

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"maps"
	"reflect"

	"orglang/go-runtime/lib/db"

	"orglang/go-runtime/adt/identity"
	"orglang/go-runtime/adt/polarity"
	"orglang/go-runtime/adt/procbind"
	"orglang/go-runtime/adt/procdec"
	"orglang/go-runtime/adt/procdef"
	"orglang/go-runtime/adt/procexp"
	"orglang/go-runtime/adt/procstep"
	"orglang/go-runtime/adt/revnum"
	"orglang/go-runtime/adt/symbol"
	"orglang/go-runtime/adt/typedef"
	"orglang/go-runtime/adt/typeexp"
	"orglang/go-runtime/adt/uniqref"
	"orglang/go-runtime/adt/uniqsym"
)

type API interface {
	Run(ExecSpec) error
	Take(procstep.StepSpec) error
	Retrieve(identity.ADT) (ExecSnap, error)
}

type ExecSpec struct {
	PoolID identity.ADT
	ExecID identity.ADT
	ProcES procexp.ExpSpec
}

type ExecRef = uniqref.ADT

type MainCfg struct {
	ExecRef ExecRef
	Bnds    map[symbol.ADT]EP2
	Acts    map[identity.ADT]procstep.StepRec
	PoolID  identity.ADT
}

// aka Configuration
type ExecSnap struct {
	ExecRef  ExecRef
	Chnls    map[symbol.ADT]EP
	StepRecs map[identity.ADT]procstep.StepRec
	PoolID   identity.ADT
	PoolRN   revnum.ADT
}

type Env struct {
	ProcDecs map[identity.ADT]procdec.DecRec
	TypeDefs map[uniqsym.ADT]typedef.DefRec
	TypeExps map[identity.ADT]typeexp.ExpRec
	Locks    map[uniqsym.ADT]Lock
}

type EP struct {
	ChnlPH symbol.ADT
	ChnlID identity.ADT
	ExpID  identity.ADT
	// provider
	PoolID identity.ADT
}

type EP2 struct {
	ChnlPH symbol.ADT
	ChnlID identity.ADT
	// provider
	PoolID identity.ADT
}

type Lock struct {
	PoolID identity.ADT
	PoolRN revnum.ADT
}

func ChnlPH(rec EP) symbol.ADT { return rec.ChnlPH }

// ответственность за процесс
type Liab struct {
	PoolID  identity.ADT
	ExecRef ExecRef
	// позитивное значение при вручении
	// негативное значение при лишении
	PoolRN revnum.ADT
}

type ExecMod struct {
	Locks []Lock
	Bnds  []procbind.BindRec
	Steps []procstep.StepRec
	Liabs []Liab
}

type MainMod struct {
	Bnds []procbind.BindRec
	Acts []procstep.StepRec
}

type service struct {
	procExecs Repo
	procDecs  procdec.Repo
	typeDefs  typedef.Repo
	typeExps  typeexp.Repo
	operator  db.Operator
	log       *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

func newService(
	procExecs Repo,
	procDecs procdec.Repo,
	typeDefs typedef.Repo,
	typeExps typeexp.Repo,
	operator db.Operator,
	l *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{procExecs, procDecs, typeDefs, typeExps, operator, l.With(name)}
}

func (s *service) Run(spec ExecSpec) (err error) {
	idAttr := slog.Any("execID", spec.ExecID)
	s.log.Debug("creation started", idAttr)
	ctx := context.Background()
	var mainCfg MainCfg
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		mainCfg, err = s.procExecs.SelectMain(ds, spec.ExecID)
		return err
	})
	if err != nil {
		s.log.Error("creation failed", idAttr)
		return err
	}
	var mainEnv Env
	err = s.checkType(spec.PoolID, mainEnv, mainCfg, spec.ProcES)
	if err != nil {
		s.log.Error("creation failed", idAttr)
		return err
	}
	mainMod, err := s.createWith(mainEnv, mainCfg, spec.ProcES)
	if err != nil {
		s.log.Error("creation failed", idAttr)
		return err
	}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.procExecs.UpdateMain(ds, mainMod)
		if err != nil {
			s.log.Error("creation failed", idAttr)
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("creation failed", idAttr)
		return err
	}
	return nil
}

func (s *service) Retrieve(execID identity.ADT) (_ ExecSnap, err error) {
	return ExecSnap{}, nil
}

func (s *service) checkType(
	poolID identity.ADT,
	mainEnv Env,
	mainCfg MainCfg,
	expSpec procexp.ExpSpec,
) error {
	chnlEP, ok := mainCfg.Bnds[expSpec.Via()]
	if !ok {
		panic("no via in main cfg")
	}
	if poolID == chnlEP.PoolID {
		return s.checkProviderMain(poolID, mainEnv, mainCfg, expSpec)
	} else {
		return s.checkClientMain(poolID, mainEnv, mainCfg, expSpec)
	}
}

func (s *service) checkProviderMain(
	poolID identity.ADT,
	mainEnv Env,
	mainCfg MainCfg,
	es procexp.ExpSpec,
) error {
	return nil
}

func (s *service) checkClientMain(
	poolID identity.ADT,
	mainEnv Env,
	mainCfg MainCfg,
	es procexp.ExpSpec,
) error {
	return nil
}

func (s *service) createWith(
	mainEnv Env,
	procCfg MainCfg,
	es procexp.ExpSpec,
) (
	procMod MainMod,
	_ error,
) {
	switch expSpec := es.(type) {
	case procexp.CallSpec:
		commChnlEP, ok := procCfg.Bnds[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("coordination failed")
			return MainMod{}, err
		}
		viaAttr := slog.Any("cordID", commChnlEP.ChnlID)
		for _, valChnlPH := range expSpec.ValChnlPHs {
			sndrValBnd := procbind.BindRec{
				ExecRef: ExecRef{ID: procCfg.ExecRef.ID, RN: -procCfg.ExecRef.RN.Next()},
				ChnlPH:  valChnlPH,
			}
			procMod.Bnds = append(procMod.Bnds, sndrValBnd)
		}
		rcvrAct := procCfg.Acts[commChnlEP.ChnlID]
		if rcvrAct == nil {
			sndrAct := procstep.MsgRec{}
			procMod.Acts = append(procMod.Acts, sndrAct)
			s.log.Debug("coordination half done", viaAttr)
			return procMod, nil
		}
		s.log.Debug("coordination succeed")
		return procMod, nil
	case procexp.SpawnSpec:
		s.log.Debug("coordination succeed")
		return procMod, nil
	default:
		panic(procexp.ErrExpTypeUnexpected(es))
	}
}

func ErrMissingChnl(want symbol.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func (s *service) Take(spec procstep.StepSpec) (err error) {
	refAttr := slog.Any("execRef", spec.ExecRef)
	s.log.Debug("taking started", refAttr)
	ctx := context.Background()
	// initial values
	poolID := spec.PoolID
	execRef := spec.ExecRef
	expSpec := spec.ProcES
	for expSpec != nil {
		var procExec ExecSnap
		err = s.operator.Implicit(ctx, func(ds db.Source) error {
			procExec, err = s.procExecs.SelectSnap(ds, execRef)
			return err
		})
		if err != nil {
			s.log.Error("taking failed", refAttr)
			return err
		}
		if len(procExec.Chnls) == 0 {
			panic("zero channels")
		}
		decIDs := procexp.CollectEnv(expSpec)
		var procDecs map[identity.ADT]procdec.DecRec
		err = s.operator.Implicit(ctx, func(ds db.Source) error {
			procDecs, err = s.procDecs.SelectEnv(ds, decIDs)
			return err
		})
		if err != nil {
			s.log.Error("taking failed", refAttr, slog.Any("decs", decIDs))
			return err
		}
		typeQNs := procdec.CollectEnv(maps.Values(procDecs))
		var typeDefs map[uniqsym.ADT]typedef.DefRec
		err = s.operator.Implicit(ctx, func(ds db.Source) error {
			typeDefs, err = s.typeDefs.SelectEnv(ds, typeQNs)
			return err
		})
		if err != nil {
			s.log.Error("taking failed", refAttr, slog.Any("types", typeQNs))
			return err
		}
		envIDs := typedef.CollectEnv(maps.Values(typeDefs))
		ctxIDs := CollectCtx(maps.Values(procExec.Chnls))
		var typeExps map[identity.ADT]typeexp.ExpRec
		err = s.operator.Implicit(ctx, func(ds db.Source) error {
			typeExps, err = s.typeExps.SelectEnv(ds, append(envIDs, ctxIDs...))
			return err
		})
		if err != nil {
			s.log.Error("taking failed", refAttr, slog.Any("env", envIDs), slog.Any("ctx", ctxIDs))
			return err
		}
		procEnv := Env{ProcDecs: procDecs, TypeDefs: typeDefs, TypeExps: typeExps}
		procCtx := convertToCtx(poolID, maps.Values(procExec.Chnls), typeExps)
		// type checking
		err = s.checkState(poolID, procEnv, procCtx, procExec, expSpec)
		if err != nil {
			s.log.Error("taking failed", refAttr)
			return err
		}
		// step taking
		nextSpec, procMod, err := s.takeWith(procEnv, procExec, expSpec)
		if err != nil {
			s.log.Error("taking failed", refAttr)
			return err
		}
		err = s.operator.Explicit(ctx, func(ds db.Source) error {
			err = s.procExecs.UpdateProc(ds, procMod)
			if err != nil {
				s.log.Error("taking failed", refAttr)
				return err
			}
			return nil
		})
		if err != nil {
			s.log.Error("taking failed", refAttr)
			return err
		}
		// next values
		poolID = nextSpec.PoolID
		execRef = nextSpec.ExecRef
		expSpec = nextSpec.ProcES
	}
	s.log.Debug("taking succeed", refAttr)
	return nil
}

func (s *service) takeWith(
	procEnv Env,
	execSnap ExecSnap,
	es procexp.ExpSpec,
) (
	stepSpec procstep.StepSpec,
	execMod ExecMod,
	_ error,
) {
	switch expSpec := es.(type) {
	case procexp.CloseSpec:
		commChnlEP, ok := execSnap.Chnls[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		viaAttr := slog.Any("chnlID", commChnlEP.ChnlID)
		sndrLock := Lock{
			PoolID: execSnap.PoolID,
			PoolRN: execSnap.PoolRN,
		}
		execMod.Locks = append(execMod.Locks, sndrLock)
		rcvrStep := execSnap.StepRecs[commChnlEP.ChnlID]
		if rcvrStep == nil {
			sndrStep := procstep.MsgRec{
				PoolID:  execSnap.PoolID,
				ExecRef: execSnap.ExecRef,
				ChnlID:  commChnlEP.ChnlID,
				PoolRN:  execSnap.PoolRN.Next(),
				ValER: procexp.CloseRec{
					CommChnlPH: expSpec.CommChnlPH,
				},
			}
			execMod.Steps = append(execMod.Steps, sndrStep)
			s.log.Debug("taking half done", viaAttr)
			return stepSpec, execMod, nil
		}
		svcStep, ok := rcvrStep.(procstep.SvcRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(rcvrStep))
		}
		switch termImpl := svcStep.ContER.(type) {
		case procexp.WaitRec:
			sndrViaBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  expSpec.CommChnlPH,
				PoolRN:  -execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, sndrViaBnd)
			rcvrViaBnd := procbind.BindRec{
				ExecRef: svcStep.ExecRef,
				ChnlPH:  termImpl.CommChnlPH,
				PoolRN:  -svcStep.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, rcvrViaBnd)
			stepSpec = procstep.StepSpec{
				PoolID:  svcStep.PoolID,
				ExecRef: svcStep.ExecRef,
				ProcES:  termImpl.ContES,
			}
			s.log.Debug("taking succeed", viaAttr)
			return stepSpec, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(svcStep.ContER))
		}
	case procexp.WaitSpec:
		commChnlEP, ok := execSnap.Chnls[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		viaAttr := slog.Any("chnlID", commChnlEP.ChnlID)
		rcvrLock := Lock{
			PoolID: execSnap.PoolID,
			PoolRN: execSnap.PoolRN,
		}
		execMod.Locks = append(execMod.Locks, rcvrLock)
		sndrStep := execSnap.StepRecs[commChnlEP.ChnlID]
		if sndrStep == nil {
			rcvrStep := procstep.SvcRec{
				PoolID:  execSnap.PoolID,
				ExecRef: execSnap.ExecRef,
				ChnlID:  commChnlEP.ChnlID,
				PoolRN:  execSnap.PoolRN.Next(),
				ContER: procexp.WaitRec{
					CommChnlPH: expSpec.CommChnlPH,
					ContES:     expSpec.ContES,
				},
			}
			execMod.Steps = append(execMod.Steps, rcvrStep)
			s.log.Debug("taking half done", viaAttr)
			return stepSpec, execMod, nil
		}
		msgStep, ok := sndrStep.(procstep.MsgRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(sndrStep))
		}
		switch termImpl := msgStep.ValER.(type) {
		case procexp.CloseRec:
			sndrViaBnd := procbind.BindRec{
				ExecRef: msgStep.ExecRef,
				ChnlPH:  termImpl.CommChnlPH,
				PoolRN:  -msgStep.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, sndrViaBnd)
			rcvrViaBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  expSpec.CommChnlPH,
				PoolRN:  -execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, rcvrViaBnd)
			stepSpec = procstep.StepSpec{
				PoolID:  execSnap.PoolID,
				ExecRef: execSnap.ExecRef,
				ProcES:  expSpec.ContES,
			}
			s.log.Debug("taking succeed", viaAttr)
			return stepSpec, execMod, nil
		case procexp.FwdRec:
			rcvrViaBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  expSpec.CommChnlPH,
				ChnlID:  termImpl.ContChnlID,
				ExpID:   commChnlEP.ExpID,
				PoolRN:  execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, rcvrViaBnd)
			stepSpec = procstep.StepSpec{
				PoolID:  execSnap.PoolID,
				ExecRef: execSnap.ExecRef,
				ProcES:  expSpec,
			}
			s.log.Debug("taking succeed", viaAttr)
			return stepSpec, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(msgStep.ValER))
		}
	case procexp.SendSpec:
		commChnlEP, ok := execSnap.Chnls[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		viaAttr := slog.Any("chnlID", commChnlEP.ChnlID)
		sndrLock := Lock{
			PoolID: execSnap.PoolID,
			PoolRN: execSnap.PoolRN,
		}
		execMod.Locks = append(execMod.Locks, sndrLock)
		viaState, ok := procEnv.TypeExps[commChnlEP.ExpID]
		if !ok {
			err := typedef.ErrMissingInEnv(commChnlEP.ExpID)
			s.log.Error("taking failed", viaAttr)
			return procstep.StepSpec{}, ExecMod{}, err
		}
		viaStateID := viaState.(typeexp.ProdRec).Next()
		valChnl, ok := execSnap.Chnls[expSpec.ValChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.ValChnlPH)
			s.log.Error("taking failed", viaAttr)
			return procstep.StepSpec{}, ExecMod{}, err
		}
		sndrValBnd := procbind.BindRec{
			ExecRef: execSnap.ExecRef,
			ChnlPH:  expSpec.ValChnlPH,
			PoolRN:  -execSnap.PoolRN.Next(),
		}
		execMod.Bnds = append(execMod.Bnds, sndrValBnd)
		rcvrStep := execSnap.StepRecs[commChnlEP.ChnlID]
		if rcvrStep == nil {
			newChnlID := identity.New()
			sndrViaBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  expSpec.CommChnlPH,
				ChnlID:  newChnlID,
				ExpID:   viaStateID,
				PoolRN:  execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, sndrViaBnd)
			sndrStep := procstep.MsgRec{
				PoolID:  execSnap.PoolID,
				ExecRef: execSnap.ExecRef,
				ChnlID:  commChnlEP.ChnlID,
				PoolRN:  execSnap.PoolRN.Next(),
				ValER: procexp.SendRec{
					CommChnlPH: expSpec.CommChnlPH,
					ContChnlID: newChnlID,
					ValChnlID:  valChnl.ChnlID,
					ValExpID:   valChnl.ExpID,
				},
			}
			execMod.Steps = append(execMod.Steps, sndrStep)
			s.log.Debug("taking half done", viaAttr)
			return stepSpec, execMod, nil
		}
		svcStep, ok := rcvrStep.(procstep.SvcRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(rcvrStep))
		}
		switch termImpl := svcStep.ContER.(type) {
		case procexp.RecvRec:
			sndrViaBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  expSpec.CommChnlPH,
				ChnlID:  termImpl.ContChnlID,
				ExpID:   viaStateID,
				PoolRN:  execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, sndrViaBnd)
			rcvrViaBnd := procbind.BindRec{
				ExecRef: svcStep.ExecRef,
				ChnlPH:  termImpl.CommChnlPH,
				ChnlID:  termImpl.ContChnlID,
				ExpID:   viaStateID,
				PoolRN:  svcStep.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, rcvrViaBnd)
			rcvrValBnd := procbind.BindRec{
				ExecRef: svcStep.ExecRef,
				ChnlPH:  termImpl.ValChnlPH,
				ChnlID:  valChnl.ChnlID,
				ExpID:   valChnl.ExpID,
				PoolRN:  svcStep.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, rcvrValBnd)
			stepSpec = procstep.StepSpec{
				PoolID:  svcStep.PoolID,
				ExecRef: svcStep.ExecRef,
				ProcES:  termImpl.ContES,
			}
			s.log.Debug("taking succeed", viaAttr)
			return stepSpec, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(svcStep.ContER))
		}
	case procexp.RecvSpec:
		commChnlEP, ok := execSnap.Chnls[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		viaAttr := slog.Any("chnlID", commChnlEP.ChnlID)
		rcvrLock := Lock{
			PoolID: execSnap.PoolID,
			PoolRN: execSnap.PoolRN,
		}
		execMod.Locks = append(execMod.Locks, rcvrLock)
		sndrSemRec := execSnap.StepRecs[commChnlEP.ChnlID]
		if sndrSemRec == nil {
			rcvrSemRec := procstep.SvcRec{
				PoolID:  execSnap.PoolID,
				ExecRef: execSnap.ExecRef,
				ChnlID:  commChnlEP.ChnlID,
				PoolRN:  execSnap.PoolRN.Next(),
				ContER: procexp.RecvRec{
					CommChnlPH: expSpec.CommChnlPH,
					ContChnlID: identity.New(),
					ValChnlPH:  expSpec.BindChnlPH,
					ContES:     expSpec.ContES,
				},
			}
			execMod.Steps = append(execMod.Steps, rcvrSemRec)
			s.log.Debug("taking half done", viaAttr)
			return stepSpec, execMod, nil
		}
		sndrMsgRec, ok := sndrSemRec.(procstep.MsgRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(sndrSemRec))
		}
		switch termRec := sndrMsgRec.ValER.(type) {
		case procexp.SendRec:
			viaState, ok := procEnv.TypeExps[commChnlEP.ExpID]
			if !ok {
				err := typedef.ErrMissingInEnv(commChnlEP.ExpID)
				s.log.Error("taking failed", viaAttr)
				return procstep.StepSpec{}, ExecMod{}, err
			}
			rcvrViaBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  expSpec.CommChnlPH,
				ChnlID:  termRec.ContChnlID,
				ExpID:   viaState.(typeexp.ProdRec).Next(),
				PoolRN:  execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, rcvrViaBnd)
			rcvrValBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  expSpec.BindChnlPH,
				ChnlID:  termRec.ValChnlID,
				ExpID:   termRec.ValExpID,
				PoolRN:  execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, rcvrValBnd)
			stepSpec = procstep.StepSpec{
				PoolID:  execSnap.PoolID,
				ExecRef: execSnap.ExecRef,
				ProcES:  expSpec.ContES,
			}
			s.log.Debug("taking succeed", viaAttr)
			return stepSpec, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(sndrMsgRec.ValER))
		}
	case procexp.LabSpec:
		commChnlEP, ok := execSnap.Chnls[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		viaAttr := slog.Any("chnlID", commChnlEP.ChnlID)
		sndrLock := Lock{
			PoolID: execSnap.PoolID,
			PoolRN: execSnap.PoolRN,
		}
		execMod.Locks = append(execMod.Locks, sndrLock)
		viaState, ok := procEnv.TypeExps[commChnlEP.ExpID]
		if !ok {
			err := typedef.ErrMissingInEnv(commChnlEP.ExpID)
			s.log.Error("taking failed", viaAttr)
			return procstep.StepSpec{}, ExecMod{}, err
		}
		viaStateID := viaState.(typeexp.SumRec).Next(expSpec.LabelQN)
		rcvrStep := execSnap.StepRecs[commChnlEP.ChnlID]
		if rcvrStep == nil {
			newViaID := identity.New()
			sndrViaBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  expSpec.CommChnlPH,
				ChnlID:  newViaID,
				ExpID:   viaStateID,
				PoolRN:  execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, sndrViaBnd)
			sndrStep := procstep.MsgRec{
				PoolID:  execSnap.PoolID,
				ExecRef: execSnap.ExecRef,
				ChnlID:  commChnlEP.ChnlID,
				PoolRN:  execSnap.PoolRN.Next(),
				ValER: procexp.LabRec{
					CommChnlPH: expSpec.CommChnlPH,
					ContChnlID: newViaID,
					LabelQN:    expSpec.LabelQN,
				},
			}
			execMod.Steps = append(execMod.Steps, sndrStep)
			s.log.Debug("taking half done", viaAttr)
			return stepSpec, execMod, nil
		}
		svcStep, ok := rcvrStep.(procstep.SvcRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(rcvrStep))
		}
		switch termImpl := svcStep.ContER.(type) {
		case procexp.CaseRec:
			sndrViaBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  expSpec.CommChnlPH,
				ChnlID:  termImpl.ContChnlID,
				ExpID:   viaStateID,
				PoolRN:  execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, sndrViaBnd)
			rcvrViaBnd := procbind.BindRec{
				ExecRef: svcStep.ExecRef,
				ChnlPH:  termImpl.CommChnlPH,
				ChnlID:  termImpl.ContChnlID,
				ExpID:   viaStateID,
				PoolRN:  svcStep.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, rcvrViaBnd)
			stepSpec = procstep.StepSpec{
				PoolID:  svcStep.PoolID,
				ExecRef: svcStep.ExecRef,
				ProcES:  termImpl.ContESs[expSpec.LabelQN],
			}
			s.log.Debug("taking succeed", viaAttr)
			return stepSpec, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(svcStep.ContER))
		}
	case procexp.CaseSpec:
		commChnlEP, ok := execSnap.Chnls[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		viaAttr := slog.Any("chnlID", commChnlEP.ChnlID)
		rcvrLock := Lock{
			PoolID: execSnap.PoolID,
			PoolRN: execSnap.PoolRN,
		}
		execMod.Locks = append(execMod.Locks, rcvrLock)
		sndrStep := execSnap.StepRecs[commChnlEP.ChnlID]
		if sndrStep == nil {
			rcvrStep := procstep.SvcRec{
				PoolID:  execSnap.PoolID,
				ExecRef: execSnap.ExecRef,
				ChnlID:  commChnlEP.ChnlID,
				PoolRN:  execSnap.PoolRN.Next(),
				ContER: procexp.CaseRec{
					CommChnlPH: expSpec.CommChnlPH,
					ContChnlID: identity.New(),
					ContESs:    expSpec.ContESs,
				},
			}
			execMod.Steps = append(execMod.Steps, rcvrStep)
			s.log.Debug("taking half done", viaAttr)
			return stepSpec, execMod, nil
		}
		msgStep, ok := sndrStep.(procstep.MsgRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(sndrStep))
		}
		switch termImpl := msgStep.ValER.(type) {
		case procexp.LabRec:
			viaState, ok := procEnv.TypeExps[commChnlEP.ExpID]
			if !ok {
				err := typedef.ErrMissingInEnv(commChnlEP.ExpID)
				s.log.Error("taking failed", viaAttr)
				return procstep.StepSpec{}, ExecMod{}, err
			}
			rcvrViaBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  expSpec.CommChnlPH,
				ChnlID:  termImpl.ContChnlID,
				ExpID:   viaState.(typeexp.SumRec).Next(termImpl.LabelQN),
				PoolRN:  execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, rcvrViaBnd)
			stepSpec = procstep.StepSpec{
				PoolID:  execSnap.PoolID,
				ExecRef: execSnap.ExecRef,
				ProcES:  expSpec.ContESs[termImpl.LabelQN],
			}
			s.log.Debug("taking succeed", viaAttr)
			return stepSpec, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(msgStep.ValER))
		}
	case procexp.SpawnSpecOld:
		rcvrSnap, ok := procEnv.Locks[expSpec.PoolQN]
		if !ok {
			err := errMissingPool(expSpec.PoolQN)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		rcvrLiab := Liab{
			ExecRef: uniqref.New(),
			PoolID:  rcvrSnap.PoolID,
			PoolRN:  rcvrSnap.PoolRN.Next(),
		}
		execMod.Liabs = append(execMod.Liabs, rcvrLiab)
		rcvrProcDec, ok := procEnv.ProcDecs[expSpec.SigID]
		if !ok {
			err := errMissingSig(expSpec.SigID)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		rcvrTypeDef, ok := procEnv.TypeDefs[rcvrProcDec.X.TypeQN]
		if !ok {
			err := errMissingRole(rcvrProcDec.X.TypeQN)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		contChnlID := identity.New()
		sndrViaBnd := procbind.BindRec{
			ExecRef: execSnap.ExecRef,
			ChnlPH:  expSpec.X,
			ChnlID:  contChnlID,
			ExpID:   rcvrTypeDef.ExpID,
			PoolRN:  execSnap.PoolRN.Next(),
		}
		execMod.Bnds = append(execMod.Bnds, sndrViaBnd)
		rcvrViaBnd := procbind.BindRec{
			ExecRef: rcvrLiab.ExecRef,
			ChnlPH:  rcvrProcDec.X.ChnlPH,
			ChnlID:  contChnlID,
			ExpID:   rcvrTypeDef.ExpID,
			PoolRN:  rcvrSnap.PoolRN.Next(),
		}
		execMod.Bnds = append(execMod.Bnds, rcvrViaBnd)
		for i, valChnlPH := range expSpec.Ys {
			valChnlEP, ok := execSnap.Chnls[valChnlPH]
			if !ok {
				err := ErrMissingChnl(valChnlPH)
				s.log.Error("taking failed")
				return procstep.StepSpec{}, ExecMod{}, err
			}
			sndrValBnd := procbind.BindRec{
				ExecRef: execSnap.ExecRef,
				ChnlPH:  valChnlPH,
				PoolRN:  -execSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, sndrValBnd)
			rcvrValBnd := procbind.BindRec{
				ExecRef: rcvrLiab.ExecRef,
				ChnlPH:  rcvrProcDec.Ys[i].ChnlPH,
				ChnlID:  valChnlEP.ChnlID,
				ExpID:   valChnlEP.ExpID,
				PoolRN:  rcvrSnap.PoolRN.Next(),
			}
			execMod.Bnds = append(execMod.Bnds, rcvrValBnd)
		}
		stepSpec = procstep.StepSpec{
			PoolID:  execSnap.PoolID,
			ExecRef: execSnap.ExecRef,
			ProcES:  expSpec.ContES,
		}
		s.log.Debug("taking succeed")
		return stepSpec, execMod, nil
	case procexp.FwdSpec:
		commChnlEP, ok := execSnap.Chnls[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		viaAttr := slog.Any("chnlID", commChnlEP.ChnlID)
		commChnlER, ok := procEnv.TypeExps[commChnlEP.ExpID]
		if !ok {
			err := typedef.ErrMissingInEnv(commChnlEP.ExpID)
			s.log.Error("taking failed", viaAttr)
			return procstep.StepSpec{}, ExecMod{}, err
		}
		valChnlEP, ok := execSnap.Chnls[expSpec.ContChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.ContChnlPH)
			s.log.Error("taking failed")
			return procstep.StepSpec{}, ExecMod{}, err
		}
		commChnlSR := execSnap.StepRecs[commChnlEP.ChnlID]
		switch commChnlER.Pol() {
		case polarity.Pos:
			switch stepRec := commChnlSR.(type) {
			case procstep.SvcRec:
				xBnd := procbind.BindRec{
					ExecRef: stepRec.ExecRef,
					ChnlPH:  stepRec.ContER.Via(),
					ChnlID:  commChnlEP.ChnlID,
					ExpID:   commChnlEP.ExpID,
					PoolRN:  stepRec.PoolRN.Next(),
				}
				execMod.Bnds = append(execMod.Bnds, xBnd)
				stepSpec = procstep.StepSpec{
					PoolID:  stepRec.PoolID,
					ExecRef: stepRec.ExecRef,
					ProcES:  stepRec.ContER,
				}
				s.log.Debug("taking succeed", viaAttr)
				return stepSpec, execMod, nil
			case procstep.MsgRec:
				yBnd := procbind.BindRec{
					ExecRef: stepRec.ExecRef,
					ChnlPH:  stepRec.ValER.Via(),
					ChnlID:  valChnlEP.ChnlID,
					ExpID:   valChnlEP.ExpID,
					PoolRN:  stepRec.PoolRN.Next(),
				}
				execMod.Bnds = append(execMod.Bnds, yBnd)
				stepSpec = procstep.StepSpec{
					PoolID:  stepRec.PoolID,
					ExecRef: stepRec.ExecRef,
					ProcES:  stepRec.ValER,
				}
				s.log.Debug("taking succeed", viaAttr)
				return stepSpec, execMod, nil
			case nil:
				xBnd := procbind.BindRec{
					ExecRef: execSnap.ExecRef,
					ChnlPH:  expSpec.CommChnlPH,
					PoolRN:  -execSnap.PoolRN.Next(),
				}
				execMod.Bnds = append(execMod.Bnds, xBnd)
				yBnd := procbind.BindRec{
					ExecRef: execSnap.ExecRef,
					ChnlPH:  expSpec.ContChnlPH,
					PoolRN:  -execSnap.PoolRN.Next(),
				}
				execMod.Bnds = append(execMod.Bnds, yBnd)
				msgStep := procstep.MsgRec{
					PoolID:  execSnap.PoolID,
					ExecRef: execSnap.ExecRef,
					ChnlID:  commChnlEP.ChnlID,
					PoolRN:  execSnap.PoolRN.Next(),
					ValER: procexp.FwdRec{
						ContChnlID: valChnlEP.ChnlID,
					},
				}
				execMod.Steps = append(execMod.Steps, msgStep)
				s.log.Debug("taking half done", viaAttr)
				return stepSpec, execMod, nil
			default:
				panic(procstep.ErrRecTypeUnexpected(commChnlSR))
			}
		case polarity.Neg:
			switch stepRec := commChnlSR.(type) {
			case procstep.SvcRec:
				yBnd := procbind.BindRec{
					ExecRef: stepRec.ExecRef,
					ChnlPH:  stepRec.ContER.Via(),
					ChnlID:  valChnlEP.ChnlID,
					ExpID:   valChnlEP.ExpID,
					PoolRN:  stepRec.PoolRN.Next(),
				}
				execMod.Bnds = append(execMod.Bnds, yBnd)
				stepSpec = procstep.StepSpec{
					PoolID:  stepRec.PoolID,
					ExecRef: stepRec.ExecRef,
					ProcES:  stepRec.ContER,
				}
				s.log.Debug("taking succeed", viaAttr)
				return stepSpec, execMod, nil
			case procstep.MsgRec:
				xBnd := procbind.BindRec{
					ExecRef: stepRec.ExecRef,
					ChnlPH:  stepRec.ValER.Via(),
					ChnlID:  commChnlEP.ChnlID,
					ExpID:   commChnlEP.ExpID,
					PoolRN:  stepRec.PoolRN.Next(),
				}
				execMod.Bnds = append(execMod.Bnds, xBnd)
				stepSpec = procstep.StepSpec{
					PoolID:  stepRec.PoolID,
					ExecRef: stepRec.ExecRef,
					ProcES:  stepRec.ValER,
				}
				s.log.Debug("taking succeed", viaAttr)
				return stepSpec, execMod, nil
			case nil:
				svcStep := procstep.SvcRec{
					PoolID:  execSnap.PoolID,
					ExecRef: execSnap.ExecRef,
					ChnlID:  commChnlEP.ChnlID,
					PoolRN:  execSnap.PoolRN.Next(),
					ContER: procexp.FwdRec{
						ContChnlID: valChnlEP.ChnlID,
					},
				}
				execMod.Steps = append(execMod.Steps, svcStep)
				s.log.Debug("taking half done", viaAttr)
				return stepSpec, execMod, nil
			default:
				panic(procstep.ErrRecTypeUnexpected(commChnlSR))
			}
		default:
			panic(typeexp.ErrPolarityUnexpected(commChnlER))
		}
	default:
		panic(procexp.ErrExpTypeUnexpected(es))
	}
}

func CollectCtx(chnls iter.Seq[EP]) []identity.ADT {
	return nil
}

func convertToCtx(poolID identity.ADT, chnlEPs iter.Seq[EP], typeExps map[identity.ADT]typeexp.ExpRec) typedef.Context {
	assets := make(map[symbol.ADT]typeexp.ExpRec, 1)
	liabs := make(map[symbol.ADT]typeexp.ExpRec, 1)
	for ep := range chnlEPs {
		if poolID == ep.PoolID {
			liabs[ep.ChnlPH] = typeExps[ep.ExpID]
		} else {
			assets[ep.ChnlPH] = typeExps[ep.ExpID]
		}
	}
	return typedef.Context{Assets: assets, Liabs: liabs}
}

func (s *service) checkState(
	poolID identity.ADT,
	procEnv Env,
	procCtx typedef.Context,
	procCfg ExecSnap,
	expSpec procexp.ExpSpec,
) error {
	chnlEP, ok := procCfg.Chnls[expSpec.Via()]
	if !ok {
		panic("no via in proc snap")
	}
	if poolID == chnlEP.PoolID {
		return s.checkProvider(poolID, procEnv, procCtx, procCfg, expSpec)
	} else {
		return s.checkClient(poolID, procEnv, procCtx, procCfg, expSpec)
	}
}

func (s *service) checkProvider(
	poolID identity.ADT,
	procEnv Env,
	procCtx typedef.Context,
	procCfg ExecSnap,
	es procexp.ExpSpec,
) error {
	switch expSpec := es.(type) {
	case procexp.CloseSpec:
		// check ctx
		if len(procCtx.Assets) > 0 {
			err := fmt.Errorf("context mismatch: want 0 items, got %v items", len(procCtx.Assets))
			s.log.Error("checking failed")
			return err
		}
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.CommChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typeexp.CheckRec(gotVia, typeexp.OneRec{})
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// no cont to check
		delete(procCtx.Liabs, expSpec.CommChnlPH)
		return nil
	case procexp.WaitSpec:
		err := procexp.ErrExpTypeMismatch(es, procexp.CloseSpec{})
		s.log.Error("checking failed")
		return err
	case procexp.SendSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.CommChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typeexp.TensorRec)
		if !ok {
			err := typeexp.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check value
		gotVal, ok := procCtx.Assets[expSpec.ValChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.ValChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typeexp.CheckRec(gotVal, wantVia.Y)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// no cont to check
		procCtx.Liabs[expSpec.CommChnlPH] = wantVia.Z
		delete(procCtx.Assets, expSpec.ValChnlPH)
		return nil
	case procexp.RecvSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.CommChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typeexp.LolliRec)
		if !ok {
			err := typeexp.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check value
		gotVal, ok := procCtx.Assets[expSpec.BindChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.BindChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typeexp.CheckRec(gotVal, wantVia.Y)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// check cont
		procCtx.Liabs[expSpec.CommChnlPH] = wantVia.Z
		procCtx.Assets[expSpec.BindChnlPH] = wantVia.Y
		return s.checkState(poolID, procEnv, procCtx, procCfg, expSpec.ContES)
	case procexp.LabSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.CommChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typeexp.PlusRec)
		if !ok {
			err := typeexp.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check label
		choice, ok := wantVia.Zs[expSpec.LabelQN]
		if !ok {
			err := fmt.Errorf("label mismatch: want %v, got %q", maps.Keys(wantVia.Zs), expSpec.LabelQN)
			s.log.Error("checking failed")
			return err
		}
		// no cont to check
		procCtx.Liabs[expSpec.CommChnlPH] = choice
		return nil
	case procexp.CaseSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.CommChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typeexp.WithRec)
		if !ok {
			err := typeexp.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check conts
		if len(expSpec.ContESs) != len(wantVia.Zs) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantVia.Zs), len(expSpec.ContESs))
			s.log.Error("checking failed")
			return err
		}
		for label, choice := range wantVia.Zs {
			cont, ok := expSpec.ContESs[label]
			if !ok {
				err := fmt.Errorf("label mismatch: want %q, got nothing", label)
				s.log.Error("checking failed")
				return err
			}
			procCtx.Liabs[expSpec.CommChnlPH] = choice
			err := s.checkState(poolID, procEnv, procCtx, procCfg, cont)
			if err != nil {
				s.log.Error("checking failed")
				return err
			}
		}
		return nil
	case procexp.FwdSpec:
		if len(procCtx.Assets) != 1 {
			err := fmt.Errorf("context mismatch: want 1 item, got %v items", len(procCtx.Assets))
			s.log.Error("checking failed")
			return err
		}
		viaSt, ok := procCtx.Liabs[expSpec.CommChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		fwdSt, ok := procCtx.Assets[expSpec.ContChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.ContChnlPH)
			s.log.Error("checking failed")
			return err
		}
		if fwdSt.Pol() != viaSt.Pol() {
			err := typeexp.ErrPolarityMismatch(fwdSt, viaSt)
			s.log.Error("checking failed")
			return err
		}
		err := typeexp.CheckRec(fwdSt, viaSt)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		delete(procCtx.Liabs, expSpec.CommChnlPH)
		delete(procCtx.Assets, expSpec.ContChnlPH)
		return nil
	default:
		panic(procexp.ErrExpTypeUnexpected(es))
	}
}

func (s *service) checkClient(
	poolID identity.ADT,
	procEnv Env,
	procCtx typedef.Context,
	procCfg ExecSnap,
	es procexp.ExpSpec,
) error {
	switch expSpec := es.(type) {
	case procexp.CloseSpec:
		err := procexp.ErrExpTypeMismatch(es, procexp.WaitSpec{})
		s.log.Error("checking failed")
		return err
	case procexp.WaitSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typeexp.OneRec)
		if !ok {
			err := typeexp.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check cont
		delete(procCtx.Assets, expSpec.CommChnlPH)
		return s.checkState(poolID, procEnv, procCtx, procCfg, expSpec.ContES)
	case procexp.SendSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typeexp.LolliRec)
		if !ok {
			err := typeexp.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check value
		gotVal, ok := procCtx.Assets[expSpec.ValChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.ValChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typeexp.CheckRec(gotVal, wantVia.Y)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		procCtx.Assets[expSpec.CommChnlPH] = wantVia.Z
		delete(procCtx.Assets, expSpec.ValChnlPH)
		return nil
	case procexp.RecvSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typeexp.TensorRec)
		if !ok {
			err := typeexp.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check value
		gotVal, ok := procCtx.Assets[expSpec.BindChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.BindChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typeexp.CheckRec(gotVal, wantVia.Y)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// check cont
		procCtx.Assets[expSpec.CommChnlPH] = wantVia.Z
		procCtx.Assets[expSpec.BindChnlPH] = wantVia.Y
		return s.checkState(poolID, procEnv, procCtx, procCfg, expSpec.ContES)
	case procexp.LabSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typeexp.WithRec)
		if !ok {
			err := typeexp.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check label
		choice, ok := wantVia.Zs[expSpec.LabelQN]
		if !ok {
			err := fmt.Errorf("label mismatch: want %v, got %q", maps.Keys(wantVia.Zs), expSpec.LabelQN)
			s.log.Error("checking failed")
			return err
		}
		procCtx.Assets[expSpec.CommChnlPH] = choice
		return nil
	case procexp.CaseSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.CommChnlPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typeexp.PlusRec)
		if !ok {
			err := typeexp.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check conts
		if len(expSpec.ContESs) != len(wantVia.Zs) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantVia.Zs), len(expSpec.ContESs))
			s.log.Error("checking failed")
			return err
		}
		for label, choice := range wantVia.Zs {
			cont, ok := expSpec.ContESs[label]
			if !ok {
				err := fmt.Errorf("label mismatch: want %q, got nothing", label)
				s.log.Error("checking failed")
				return err
			}
			procCtx.Assets[expSpec.CommChnlPH] = choice
			err := s.checkState(poolID, procEnv, procCtx, procCfg, cont)
			if err != nil {
				s.log.Error("checking failed")
				return err
			}
		}
		return nil
	case procexp.SpawnSpecOld:
		procSig, ok := procEnv.ProcDecs[expSpec.SigID]
		if !ok {
			err := procdec.ErrRootMissingInEnv(expSpec.SigID)
			s.log.Error("checking failed")
			return err
		}
		// check vals
		if len(expSpec.Ys) != len(procSig.Ys) {
			err := fmt.Errorf("context mismatch: want %v items, got %v items", len(procSig.Ys), len(expSpec.Ys))
			s.log.Error("checking failed", slog.Any("want", procSig.Ys), slog.Any("got", expSpec.Ys))
			return err
		}
		if len(expSpec.Ys) == 0 {
			return nil
		}
		for i, ep := range procSig.Ys {
			valRole, ok := procEnv.TypeDefs[ep.TypeQN]
			if !ok {
				err := typedef.ErrSymMissingInEnv(ep.TypeQN)
				s.log.Error("checking failed")
				return err
			}
			wantVal, ok := procEnv.TypeExps[valRole.ExpID]
			if !ok {
				err := typedef.ErrMissingInEnv(valRole.ExpID)
				s.log.Error("checking failed")
				return err
			}
			gotVal, ok := procCtx.Assets[expSpec.Ys[i]]
			if !ok {
				err := procdef.ErrMissingInCtx(ep.ChnlPH)
				s.log.Error("checking failed")
				return err
			}
			err := typeexp.CheckRec(gotVal, wantVal)
			if err != nil {
				s.log.Error("checking failed", slog.Any("want", wantVal), slog.Any("got", gotVal))
				return err
			}
			delete(procCtx.Assets, expSpec.Ys[i])
		}
		// check via
		viaRole, ok := procEnv.TypeDefs[procSig.X.TypeQN]
		if !ok {
			err := typedef.ErrSymMissingInEnv(procSig.X.TypeQN)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := procEnv.TypeExps[viaRole.ExpID]
		if !ok {
			err := typedef.ErrMissingInEnv(viaRole.ExpID)
			s.log.Error("checking failed")
			return err
		}
		// check cont
		procCtx.Assets[expSpec.X] = wantVia
		return s.checkState(poolID, procEnv, procCtx, procCfg, expSpec.ContES)
	default:
		panic(procexp.ErrExpTypeUnexpected(es))
	}
}

func errOptimisticUpdate(got revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}

func errMissingPool(want uniqsym.ADT) error {
	return fmt.Errorf("pool missing in env: %v", want)
}

func errMissingSig(want identity.ADT) error {
	return fmt.Errorf("sig missing in env: %v", want)
}

func errMissingRole(want uniqsym.ADT) error {
	return fmt.Errorf("role missing in env: %v", want)
}
