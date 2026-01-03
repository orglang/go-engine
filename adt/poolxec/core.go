package poolxec

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/exp/maps"

	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/polarity"
	"orglang/orglang/adt/qualsym"
	"orglang/orglang/adt/revnum"

	"orglang/orglang/adt/pooldef"
	"orglang/orglang/adt/procdec"
	"orglang/orglang/adt/procdef"
	"orglang/orglang/adt/procxec"
	"orglang/orglang/adt/typedef"
)

// Port
type API interface {
	Create(ExecSpec) (ExecRef, error)
	Retrieve(identity.ADT) (ExecSnap, error)
	RetreiveRefs() ([]ExecRef, error)
	Spawn(procxec.ExecSpec) (procxec.ExecRef, error)
	Take(StepSpec) error
	Poll(PollSpec) (procxec.ExecRef, error)
}

type ExecSpec struct {
	PoolQN qualsym.ADT
	SupID  identity.ADT
}

type ExecRef struct {
	ExecID identity.ADT
	ProcID identity.ADT // main
}

type ExecRec struct {
	ExecID identity.ADT
	ProcID identity.ADT // main
	SupID  identity.ADT
	PoolRN revnum.ADT
}

type ExecSnap struct {
	ExecID identity.ADT
	Title  string
	Subs   []ExecRef
}

type StepSpec struct {
	PoolID identity.ADT
	ProcID identity.ADT
	ProcTS procdef.TermSpec
}

type PollSpec struct {
	PoolID identity.ADT
	PoolTS pooldef.TermSpec
}

type service struct {
	pools    Repo
	procs    procdec.Repo
	types    typedef.Repo
	operator sd.Operator
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

func newService(
	pools Repo,
	procs procdec.Repo,
	types typedef.Repo,
	operator sd.Operator,
	l *slog.Logger,
) *service {
	return &service{pools, procs, types, operator, l}
}

func (s *service) Create(spec ExecSpec) (ExecRef, error) {
	ctx := context.Background()
	s.log.Debug("creation started", slog.Any("spec", spec))
	impl := ExecRec{
		ExecID: identity.New(),
		ProcID: identity.New(),
		SupID:  spec.SupID,
		PoolRN: revnum.Initial(),
	}
	liab := procxec.Liab{
		PoolID: impl.ExecID,
		ProcID: impl.ProcID,
		PoolRN: impl.PoolRN,
	}
	err := s.operator.Explicit(ctx, func(ds sd.Source) error {
		err := s.pools.Insert(ds, impl)
		if err != nil {
			s.log.Error("creation failed")
			return err
		}
		err = s.pools.InsertLiab(ds, liab)
		if err != nil {
			s.log.Error("creation failed")
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("creation failed")
		return ExecRef{}, err
	}
	s.log.Debug("creation succeed", slog.Any("poolID", impl.ExecID))
	return ConvertRecToRef(impl), nil
}

func (s *service) Poll(spec PollSpec) (procxec.ExecRef, error) {
	switch spec.PoolTS.(type) {
	default:
		return procxec.ExecRef{}, nil
	}
	return procxec.ExecRef{}, nil
}

func (s *service) Spawn(spec procxec.ExecSpec) (_ procxec.ExecRef, err error) {
	procAttr := slog.Any("procID", spec.ExecID)
	s.log.Debug("spawning started", procAttr)
	return procxec.ExecRef{}, nil
}

func (s *service) Take(spec StepSpec) (err error) {
	idAttr := slog.Any("procID", spec.ProcID)
	s.log.Debug("taking started", idAttr)
	ctx := context.Background()
	// initial values
	poolID := spec.PoolID
	procID := spec.ProcID
	termSpec := spec.ProcTS
	for termSpec != nil {
		var procCfg procxec.Cfg
		err = s.operator.Implicit(ctx, func(ds sd.Source) error {
			procCfg, err = s.pools.SelectProc(ds, procID)
			return err
		})
		if err != nil {
			s.log.Error("taking failed", idAttr)
			return err
		}
		if len(procCfg.Chnls) == 0 {
			panic("zero channels")
		}
		sigIDs := procdef.CollectEnv(termSpec)
		var sigs map[identity.ADT]procdec.DecRec
		err = s.operator.Implicit(ctx, func(ds sd.Source) error {
			sigs, err = s.procs.SelectEnv(ds, sigIDs)
			return err
		})
		if err != nil {
			s.log.Error("taking failed", idAttr, slog.Any("sigs", sigIDs))
			return err
		}
		typeQNs := procdec.CollectEnv(maps.Values(sigs))
		var types map[qualsym.ADT]typedef.DefRec
		err = s.operator.Implicit(ctx, func(ds sd.Source) error {
			types, err = s.types.SelectTypeEnv(ds, typeQNs)
			return err
		})
		if err != nil {
			s.log.Error("taking failed", idAttr, slog.Any("types", typeQNs))
			return err
		}
		envIDs := typedef.CollectEnv(maps.Values(types))
		ctxIDs := CollectCtx(maps.Values(procCfg.Chnls))
		var terms map[identity.ADT]typedef.TermRec
		err = s.operator.Implicit(ctx, func(ds sd.Source) error {
			terms, err = s.types.SelectTermEnv(ds, append(envIDs, ctxIDs...))
			return err
		})
		if err != nil {
			s.log.Error("taking failed", idAttr, slog.Any("env", envIDs), slog.Any("ctx", ctxIDs))
			return err
		}
		procEnv := procxec.Env{ProcDecs: sigs, TypeDefs: types, TypeTerms: terms}
		procCtx := convertToCtx(poolID, maps.Values(procCfg.Chnls), terms)
		// type checking
		err = s.checkState(poolID, procEnv, procCtx, procCfg, termSpec)
		if err != nil {
			s.log.Error("taking failed", idAttr)
			return err
		}
		// step taking
		nextSpec, procMod, err := s.takeWith(procEnv, procCfg, termSpec)
		if err != nil {
			s.log.Error("taking failed", idAttr)
			return err
		}
		err = s.operator.Explicit(ctx, func(ds sd.Source) error {
			err = s.pools.UpdateProc(ds, procMod)
			if err != nil {
				s.log.Error("taking failed", idAttr)
				return err
			}
			return nil
		})
		if err != nil {
			s.log.Error("taking failed", idAttr)
			return err
		}
		// next values
		poolID = nextSpec.PoolID
		procID = nextSpec.ProcID
		termSpec = nextSpec.ProcTS
	}
	s.log.Debug("taking succeed", idAttr)
	return nil
}

func (s *service) takeWith(
	procEnv procxec.Env,
	procCfg procxec.Cfg,
	ts procdef.TermSpec,
) (
	tranSpec StepSpec,
	procMod procxec.Mod,
	_ error,
) {
	switch termSpec := ts.(type) {
	case procdef.CloseSpec:
		viaChnl, ok := procCfg.Chnls[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCfg(termSpec.CommPH)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		sndrLock := procxec.Lock{
			PoolID: procCfg.PoolID,
			PoolRN: procCfg.PoolRN,
		}
		procMod.Locks = append(procMod.Locks, sndrLock)
		rcvrStep := procCfg.Steps[viaChnl.ChnlID]
		if rcvrStep == nil {
			sndrStep := procxec.MsgRec{
				PoolID: procCfg.PoolID,
				ProcID: procCfg.ExecID,
				ChnlID: viaChnl.ChnlID,
				PoolRN: procCfg.PoolRN.Next(),
				Val: procdef.CloseRec{
					X: termSpec.CommPH,
				},
			}
			procMod.Steps = append(procMod.Steps, sndrStep)
			s.log.Debug("taking half done", viaAttr)
			return tranSpec, procMod, nil
		}
		svcStep, ok := rcvrStep.(procxec.SvcRec)
		if !ok {
			panic(procxec.ErrRootTypeUnexpected(rcvrStep))
		}
		switch termImpl := svcStep.Cont.(type) {
		case procdef.WaitRec:
			sndrViaBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: termSpec.CommPH,
				PoolRN: -procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, sndrViaBnd)
			rcvrViaBnd := procxec.Bnd{
				ProcID: svcStep.ProcID,
				ChnlPH: termImpl.X,
				PoolRN: -svcStep.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, rcvrViaBnd)
			tranSpec = StepSpec{
				PoolID: svcStep.PoolID,
				ProcID: svcStep.ProcID,
				ProcTS: termImpl.ContTS,
			}
			s.log.Debug("taking succeed", viaAttr)
			return tranSpec, procMod, nil
		default:
			panic(procdef.ErrRecTypeUnexpected(svcStep.Cont))
		}
	case procdef.WaitSpec:
		viaChnl, ok := procCfg.Chnls[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCfg(termSpec.CommPH)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		rcvrLock := procxec.Lock{
			PoolID: procCfg.PoolID,
			PoolRN: procCfg.PoolRN,
		}
		procMod.Locks = append(procMod.Locks, rcvrLock)
		sndrStep := procCfg.Steps[viaChnl.ChnlID]
		if sndrStep == nil {
			rcvrStep := procxec.SvcRec{
				PoolID: procCfg.PoolID,
				ProcID: procCfg.ExecID,
				ChnlID: viaChnl.ChnlID,
				PoolRN: procCfg.PoolRN.Next(),
				Cont: procdef.WaitRec{
					X:      termSpec.CommPH,
					ContTS: termSpec.ContTS,
				},
			}
			procMod.Steps = append(procMod.Steps, rcvrStep)
			s.log.Debug("taking half done", viaAttr)
			return tranSpec, procMod, nil
		}
		msgStep, ok := sndrStep.(procxec.MsgRec)
		if !ok {
			panic(procxec.ErrRootTypeUnexpected(sndrStep))
		}
		switch termImpl := msgStep.Val.(type) {
		case procdef.CloseRec:
			sndrViaBnd := procxec.Bnd{
				ProcID: msgStep.ProcID,
				ChnlPH: termImpl.X,
				PoolRN: -msgStep.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, sndrViaBnd)
			rcvrViaBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: termSpec.CommPH,
				PoolRN: -procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, rcvrViaBnd)
			tranSpec = StepSpec{
				PoolID: procCfg.PoolID,
				ProcID: procCfg.ExecID,
				ProcTS: termSpec.ContTS,
			}
			s.log.Debug("taking succeed", viaAttr)
			return tranSpec, procMod, nil
		case procdef.FwdRec:
			rcvrViaBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: termSpec.CommPH,
				ChnlID: termImpl.B,
				TermID: viaChnl.TermID,
				PoolRN: procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, rcvrViaBnd)
			tranSpec = StepSpec{
				PoolID: procCfg.PoolID,
				ProcID: procCfg.ExecID,
				ProcTS: termSpec,
			}
			s.log.Debug("taking succeed", viaAttr)
			return tranSpec, procMod, nil
		default:
			panic(procdef.ErrRecTypeUnexpected(msgStep.Val))
		}
	case procdef.SendSpec:
		viaChnl, ok := procCfg.Chnls[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCfg(termSpec.CommPH)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		sndrLock := procxec.Lock{
			PoolID: procCfg.PoolID,
			PoolRN: procCfg.PoolRN,
		}
		procMod.Locks = append(procMod.Locks, sndrLock)
		viaState, ok := procEnv.TypeTerms[viaChnl.TermID]
		if !ok {
			err := typedef.ErrMissingInEnv(viaChnl.TermID)
			s.log.Error("taking failed", viaAttr)
			return StepSpec{}, procxec.Mod{}, err
		}
		viaStateID := viaState.(typedef.ProdRec).Next()
		valChnl, ok := procCfg.Chnls[termSpec.ValPH]
		if !ok {
			err := procdef.ErrMissingInCfg(termSpec.ValPH)
			s.log.Error("taking failed", viaAttr)
			return StepSpec{}, procxec.Mod{}, err
		}
		sndrValBnd := procxec.Bnd{
			ProcID: procCfg.ExecID,
			ChnlPH: termSpec.ValPH,
			PoolRN: -procCfg.PoolRN.Next(),
		}
		procMod.Bnds = append(procMod.Bnds, sndrValBnd)
		rcvrStep := procCfg.Steps[viaChnl.ChnlID]
		if rcvrStep == nil {
			newChnlID := identity.New()
			sndrViaBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: termSpec.CommPH,
				ChnlID: newChnlID,
				TermID: viaStateID,
				PoolRN: procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, sndrViaBnd)
			sndrStep := procxec.MsgRec{
				PoolID: procCfg.PoolID,
				ProcID: procCfg.ExecID,
				ChnlID: viaChnl.ChnlID,
				PoolRN: procCfg.PoolRN.Next(),
				Val: procdef.SendRec{
					X:      termSpec.CommPH,
					A:      newChnlID,
					B:      valChnl.ChnlID,
					TermID: valChnl.TermID,
				},
			}
			procMod.Steps = append(procMod.Steps, sndrStep)
			s.log.Debug("taking half done", viaAttr)
			return tranSpec, procMod, nil
		}
		svcStep, ok := rcvrStep.(procxec.SvcRec)
		if !ok {
			panic(procxec.ErrRootTypeUnexpected(rcvrStep))
		}
		switch termImpl := svcStep.Cont.(type) {
		case procdef.RecvRec:
			sndrViaBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: termSpec.CommPH,
				ChnlID: termImpl.A,
				TermID: viaStateID,
				PoolRN: procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, sndrViaBnd)
			rcvrViaBnd := procxec.Bnd{
				ProcID: svcStep.ProcID,
				ChnlPH: termImpl.X,
				ChnlID: termImpl.A,
				TermID: viaStateID,
				PoolRN: svcStep.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, rcvrViaBnd)
			rcvrValBnd := procxec.Bnd{
				ProcID: svcStep.ProcID,
				ChnlPH: termImpl.Y,
				ChnlID: valChnl.ChnlID,
				TermID: valChnl.TermID,
				PoolRN: svcStep.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, rcvrValBnd)
			tranSpec = StepSpec{
				PoolID: svcStep.PoolID,
				ProcID: svcStep.ProcID,
				ProcTS: termImpl.ContTS,
			}
			s.log.Debug("taking succeed", viaAttr)
			return tranSpec, procMod, nil
		default:
			panic(procdef.ErrRecTypeUnexpected(svcStep.Cont))
		}
	case procdef.RecvSpec:
		viaChnl, ok := procCfg.Chnls[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCfg(termSpec.CommPH)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		rcvrLock := procxec.Lock{
			PoolID: procCfg.PoolID,
			PoolRN: procCfg.PoolRN,
		}
		procMod.Locks = append(procMod.Locks, rcvrLock)
		sndrSemRec := procCfg.Steps[viaChnl.ChnlID]
		if sndrSemRec == nil {
			rcvrSemRec := procxec.SvcRec{
				PoolID: procCfg.PoolID,
				ProcID: procCfg.ExecID,
				ChnlID: viaChnl.ChnlID,
				PoolRN: procCfg.PoolRN.Next(),
				Cont: procdef.RecvRec{
					X:      termSpec.CommPH,
					A:      identity.New(),
					Y:      termSpec.BindPH,
					ContTS: termSpec.ContTS,
				},
			}
			procMod.Steps = append(procMod.Steps, rcvrSemRec)
			s.log.Debug("taking half done", viaAttr)
			return tranSpec, procMod, nil
		}
		sndrMsgRec, ok := sndrSemRec.(procxec.MsgRec)
		if !ok {
			panic(procxec.ErrRootTypeUnexpected(sndrSemRec))
		}
		switch termRec := sndrMsgRec.Val.(type) {
		case procdef.SendRec:
			viaState, ok := procEnv.TypeTerms[viaChnl.TermID]
			if !ok {
				err := typedef.ErrMissingInEnv(viaChnl.TermID)
				s.log.Error("taking failed", viaAttr)
				return StepSpec{}, procxec.Mod{}, err
			}
			rcvrViaBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: termSpec.CommPH,
				ChnlID: termRec.A,
				TermID: viaState.(typedef.ProdRec).Next(),
				PoolRN: procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, rcvrViaBnd)
			rcvrValBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: termSpec.BindPH,
				ChnlID: termRec.B,
				TermID: termRec.TermID,
				PoolRN: procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, rcvrValBnd)
			tranSpec = StepSpec{
				PoolID: procCfg.PoolID,
				ProcID: procCfg.ExecID,
				ProcTS: termSpec.ContTS,
			}
			s.log.Debug("taking succeed", viaAttr)
			return tranSpec, procMod, nil
		default:
			panic(procdef.ErrRecTypeUnexpected(sndrMsgRec.Val))
		}
	case procdef.LabSpec:
		viaChnl, ok := procCfg.Chnls[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCfg(termSpec.CommPH)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		sndrLock := procxec.Lock{
			PoolID: procCfg.PoolID,
			PoolRN: procCfg.PoolRN,
		}
		procMod.Locks = append(procMod.Locks, sndrLock)
		viaState, ok := procEnv.TypeTerms[viaChnl.TermID]
		if !ok {
			err := typedef.ErrMissingInEnv(viaChnl.TermID)
			s.log.Error("taking failed", viaAttr)
			return StepSpec{}, procxec.Mod{}, err
		}
		viaStateID := viaState.(typedef.SumRec).Next(termSpec.Label)
		rcvrStep := procCfg.Steps[viaChnl.ChnlID]
		if rcvrStep == nil {
			newViaID := identity.New()
			sndrViaBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: termSpec.CommPH,
				ChnlID: newViaID,
				TermID: viaStateID,
				PoolRN: procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, sndrViaBnd)
			sndrStep := procxec.MsgRec{
				PoolID: procCfg.PoolID,
				ProcID: procCfg.ExecID,
				ChnlID: viaChnl.ChnlID,
				PoolRN: procCfg.PoolRN.Next(),
				Val: procdef.LabRec{
					X:     termSpec.CommPH,
					A:     newViaID,
					Label: termSpec.Label,
				},
			}
			procMod.Steps = append(procMod.Steps, sndrStep)
			s.log.Debug("taking half done", viaAttr)
			return tranSpec, procMod, nil
		}
		svcStep, ok := rcvrStep.(procxec.SvcRec)
		if !ok {
			panic(procxec.ErrRootTypeUnexpected(rcvrStep))
		}
		switch termImpl := svcStep.Cont.(type) {
		case procdef.CaseRec:
			sndrViaBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: termSpec.CommPH,
				ChnlID: termImpl.A,
				TermID: viaStateID,
				PoolRN: procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, sndrViaBnd)
			rcvrViaBnd := procxec.Bnd{
				ProcID: svcStep.ProcID,
				ChnlPH: termImpl.X,
				ChnlID: termImpl.A,
				TermID: viaStateID,
				PoolRN: svcStep.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, rcvrViaBnd)
			tranSpec = StepSpec{
				PoolID: svcStep.PoolID,
				ProcID: svcStep.ProcID,
				ProcTS: termImpl.ContTSs[termSpec.Label],
			}
			s.log.Debug("taking succeed", viaAttr)
			return tranSpec, procMod, nil
		default:
			panic(procdef.ErrRecTypeUnexpected(svcStep.Cont))
		}
	case procdef.CaseSpec:
		viaChnl, ok := procCfg.Chnls[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCfg(termSpec.CommPH)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		rcvrLock := procxec.Lock{
			PoolID: procCfg.PoolID,
			PoolRN: procCfg.PoolRN,
		}
		procMod.Locks = append(procMod.Locks, rcvrLock)
		sndrStep := procCfg.Steps[viaChnl.ChnlID]
		if sndrStep == nil {
			rcvrStep := procxec.SvcRec{
				PoolID: procCfg.PoolID,
				ProcID: procCfg.ExecID,
				ChnlID: viaChnl.ChnlID,
				PoolRN: procCfg.PoolRN.Next(),
				Cont: procdef.CaseRec{
					X:       termSpec.CommPH,
					A:       identity.New(),
					ContTSs: termSpec.ContTSs,
				},
			}
			procMod.Steps = append(procMod.Steps, rcvrStep)
			s.log.Debug("taking half done", viaAttr)
			return tranSpec, procMod, nil
		}
		msgStep, ok := sndrStep.(procxec.MsgRec)
		if !ok {
			panic(procxec.ErrRootTypeUnexpected(sndrStep))
		}
		switch termImpl := msgStep.Val.(type) {
		case procdef.LabRec:
			viaState, ok := procEnv.TypeTerms[viaChnl.TermID]
			if !ok {
				err := typedef.ErrMissingInEnv(viaChnl.TermID)
				s.log.Error("taking failed", viaAttr)
				return StepSpec{}, procxec.Mod{}, err
			}
			rcvrViaBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: termSpec.CommPH,
				ChnlID: termImpl.A,
				TermID: viaState.(typedef.SumRec).Next(termImpl.Label),
				PoolRN: procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, rcvrViaBnd)
			tranSpec = StepSpec{
				PoolID: procCfg.PoolID,
				ProcID: procCfg.ExecID,
				ProcTS: termSpec.ContTSs[termImpl.Label],
			}
			s.log.Debug("taking succeed", viaAttr)
			return tranSpec, procMod, nil
		default:
			panic(procdef.ErrRecTypeUnexpected(msgStep.Val))
		}
	case procdef.SpawnSpecOld:
		rcvrSnap, ok := procEnv.Locks[termSpec.PoolQN]
		if !ok {
			err := errMissingPool(termSpec.PoolQN)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		rcvrLiab := procxec.Liab{
			ProcID: identity.New(),
			PoolID: rcvrSnap.PoolID,
			PoolRN: rcvrSnap.PoolRN.Next(),
		}
		procMod.Liabs = append(procMod.Liabs, rcvrLiab)
		rcvrSig, ok := procEnv.ProcDecs[termSpec.SigID]
		if !ok {
			err := errMissingSig(termSpec.SigID)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		rcvrRole, ok := procEnv.TypeDefs[rcvrSig.X.TypeQN]
		if !ok {
			err := errMissingRole(rcvrSig.X.TypeQN)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		newViaID := identity.New()
		sndrViaBnd := procxec.Bnd{
			ProcID: procCfg.ExecID,
			ChnlPH: termSpec.X,
			ChnlID: newViaID,
			TermID: rcvrRole.TermID,
			PoolRN: procCfg.PoolRN.Next(),
		}
		procMod.Bnds = append(procMod.Bnds, sndrViaBnd)
		rcvrViaBnd := procxec.Bnd{
			ProcID: rcvrLiab.ProcID,
			ChnlPH: rcvrSig.X.BindPH,
			ChnlID: newViaID,
			TermID: rcvrRole.TermID,
			PoolRN: rcvrSnap.PoolRN.Next(),
		}
		procMod.Bnds = append(procMod.Bnds, rcvrViaBnd)
		for i, chnlPH := range termSpec.Ys {
			valChnl, ok := procCfg.Chnls[chnlPH]
			if !ok {
				err := procxec.ErrMissingChnl(chnlPH)
				s.log.Error("taking failed")
				return StepSpec{}, procxec.Mod{}, err
			}
			sndrValBnd := procxec.Bnd{
				ProcID: procCfg.ExecID,
				ChnlPH: chnlPH,
				PoolRN: -procCfg.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, sndrValBnd)
			rcvrValBnd := procxec.Bnd{
				ProcID: rcvrLiab.ProcID,
				ChnlPH: rcvrSig.Ys[i].BindPH,
				ChnlID: valChnl.ChnlID,
				TermID: valChnl.TermID,
				PoolRN: rcvrSnap.PoolRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, rcvrValBnd)
		}
		tranSpec = StepSpec{
			PoolID: procCfg.PoolID,
			ProcID: procCfg.ExecID,
			ProcTS: termSpec.ContTS,
		}
		s.log.Debug("taking succeed")
		return tranSpec, procMod, nil
	case procdef.FwdSpec:
		viaChnl, ok := procCfg.Chnls[termSpec.X]
		if !ok {
			err := procdef.ErrMissingInCfg(termSpec.X)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		viaState, ok := procEnv.TypeTerms[viaChnl.TermID]
		if !ok {
			err := typedef.ErrMissingInEnv(viaChnl.TermID)
			s.log.Error("taking failed", viaAttr)
			return StepSpec{}, procxec.Mod{}, err
		}
		valChnl, ok := procCfg.Chnls[termSpec.Y]
		if !ok {
			err := procdef.ErrMissingInCfg(termSpec.Y)
			s.log.Error("taking failed")
			return StepSpec{}, procxec.Mod{}, err
		}
		vs := procCfg.Steps[viaChnl.ChnlID]
		switch viaState.Pol() {
		case polarity.Pos:
			switch viaStep := vs.(type) {
			case procxec.SvcRec:
				xBnd := procxec.Bnd{
					ProcID: viaStep.ProcID,
					ChnlPH: viaStep.Cont.Via(),
					ChnlID: viaChnl.ChnlID,
					TermID: viaChnl.TermID,
					PoolRN: viaStep.PoolRN.Next(),
				}
				procMod.Bnds = append(procMod.Bnds, xBnd)
				tranSpec = StepSpec{
					PoolID: viaStep.PoolID,
					ProcID: viaStep.ProcID,
					ProcTS: viaStep.Cont,
				}
				s.log.Debug("taking succeed", viaAttr)
				return tranSpec, procMod, nil
			case procxec.MsgRec:
				yBnd := procxec.Bnd{
					ProcID: viaStep.ProcID,
					ChnlPH: viaStep.Val.Via(),
					ChnlID: valChnl.ChnlID,
					TermID: valChnl.TermID,
					PoolRN: viaStep.PoolRN.Next(),
				}
				procMod.Bnds = append(procMod.Bnds, yBnd)
				tranSpec = StepSpec{
					PoolID: viaStep.PoolID,
					ProcID: viaStep.ProcID,
					ProcTS: viaStep.Val,
				}
				s.log.Debug("taking succeed", viaAttr)
				return tranSpec, procMod, nil
			case nil:
				xBnd := procxec.Bnd{
					ProcID: procCfg.ExecID,
					ChnlPH: termSpec.X,
					PoolRN: -procCfg.PoolRN.Next(),
				}
				procMod.Bnds = append(procMod.Bnds, xBnd)
				yBnd := procxec.Bnd{
					ProcID: procCfg.ExecID,
					ChnlPH: termSpec.Y,
					PoolRN: -procCfg.PoolRN.Next(),
				}
				procMod.Bnds = append(procMod.Bnds, yBnd)
				msgStep := procxec.MsgRec{
					PoolID: procCfg.PoolID,
					ProcID: procCfg.ExecID,
					ChnlID: viaChnl.ChnlID,
					PoolRN: procCfg.PoolRN.Next(),
					Val: procdef.FwdRec{
						B: valChnl.ChnlID,
					},
				}
				procMod.Steps = append(procMod.Steps, msgStep)
				s.log.Debug("taking half done", viaAttr)
				return tranSpec, procMod, nil
			default:
				panic(procxec.ErrRootTypeUnexpected(vs))
			}
		case polarity.Neg:
			switch viaStep := vs.(type) {
			case procxec.SvcRec:
				yBnd := procxec.Bnd{
					ProcID: viaStep.ProcID,
					ChnlPH: viaStep.Cont.Via(),
					ChnlID: valChnl.ChnlID,
					TermID: valChnl.TermID,
					PoolRN: viaStep.PoolRN.Next(),
				}
				procMod.Bnds = append(procMod.Bnds, yBnd)
				tranSpec = StepSpec{
					PoolID: viaStep.PoolID,
					ProcID: viaStep.ProcID,
					ProcTS: viaStep.Cont,
				}
				s.log.Debug("taking succeed", viaAttr)
				return tranSpec, procMod, nil
			case procxec.MsgRec:
				xBnd := procxec.Bnd{
					ProcID: viaStep.ProcID,
					ChnlPH: viaStep.Val.Via(),
					ChnlID: viaChnl.ChnlID,
					TermID: viaChnl.TermID,
					PoolRN: viaStep.PoolRN.Next(),
				}
				procMod.Bnds = append(procMod.Bnds, xBnd)
				tranSpec = StepSpec{
					PoolID: viaStep.PoolID,
					ProcID: viaStep.ProcID,
					ProcTS: viaStep.Val,
				}
				s.log.Debug("taking succeed", viaAttr)
				return tranSpec, procMod, nil
			case nil:
				svcStep := procxec.SvcRec{
					PoolID: procCfg.PoolID,
					ProcID: procCfg.ExecID,
					ChnlID: viaChnl.ChnlID,
					PoolRN: procCfg.PoolRN.Next(),
					Cont: procdef.FwdRec{
						B: valChnl.ChnlID,
					},
				}
				procMod.Steps = append(procMod.Steps, svcStep)
				s.log.Debug("taking half done", viaAttr)
				return tranSpec, procMod, nil
			default:
				panic(procxec.ErrRootTypeUnexpected(vs))
			}
		default:
			panic(typedef.ErrPolarityUnexpected(viaState))
		}
	default:
		panic(procdef.ErrTermTypeUnexpected(ts))
	}
}

func (s *service) Retrieve(poolID identity.ADT) (snap ExecSnap, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds sd.Source) error {
		snap, err = s.pools.SelectSubs(ds, poolID)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("poolID", poolID))
		return ExecSnap{}, err
	}
	return snap, nil
}

func (s *service) RetreiveRefs() (refs []ExecRef, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds sd.Source) error {
		refs, err = s.pools.SelectRefs(ds)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}

func CollectCtx(chnls []procxec.EP) []identity.ADT {
	return nil
}

func convertToCtx(poolID identity.ADT, chnls []procxec.EP, types map[identity.ADT]typedef.TermRec) typedef.Context {
	assets := make(map[qualsym.ADT]typedef.TermRec, len(chnls)-1)
	liabs := make(map[qualsym.ADT]typedef.TermRec, 1)
	for _, ch := range chnls {
		if poolID == ch.PoolID {
			liabs[ch.ChnlPH] = types[ch.TermID]
		} else {
			assets[ch.ChnlPH] = types[ch.TermID]
		}
	}
	return typedef.Context{Assets: assets, Liabs: liabs}
}

func (s *service) checkState(
	poolID identity.ADT,
	procEnv procxec.Env,
	procCtx typedef.Context,
	procCfg procxec.Cfg,
	termSpec procdef.TermSpec,
) error {
	ch, ok := procCfg.Chnls[termSpec.Via()]
	if !ok {
		panic("no via in proc snap")
	}
	if poolID == ch.PoolID {
		return s.checkProvider(poolID, procEnv, procCtx, procCfg, termSpec)
	} else {
		return s.checkClient(poolID, procEnv, procCtx, procCfg, termSpec)
	}
}

func (s *service) checkProvider(
	poolID identity.ADT,
	procEnv procxec.Env,
	procCtx typedef.Context,
	procCfg procxec.Cfg,
	ts procdef.TermSpec,
) error {
	switch termSpec := ts.(type) {
	case procdef.CloseSpec:
		// check ctx
		if len(procCtx.Assets) > 0 {
			err := fmt.Errorf("context mismatch: want 0 items, got %v items", len(procCtx.Assets))
			s.log.Error("checking failed")
			return err
		}
		// check via
		gotVia, ok := procCtx.Liabs[termSpec.CommPH]
		if !ok {
			err := typedef.ErrMissingInCtx(termSpec.CommPH)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(gotVia, typedef.OneRec{})
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// no cont to check
		delete(procCtx.Liabs, termSpec.CommPH)
		return nil
	case procdef.WaitSpec:
		err := procdef.ErrTermTypeMismatch(ts, procdef.CloseSpec{})
		s.log.Error("checking failed")
		return err
	case procdef.SendSpec:
		// check via
		gotVia, ok := procCtx.Liabs[termSpec.CommPH]
		if !ok {
			err := typedef.ErrMissingInCtx(termSpec.CommPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typedef.TensorRec)
		if !ok {
			err := typedef.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check value
		gotVal, ok := procCtx.Assets[termSpec.ValPH]
		if !ok {
			err := procdef.ErrMissingInCtx(termSpec.ValPH)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(gotVal, wantVia.Y)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// no cont to check
		procCtx.Liabs[termSpec.CommPH] = wantVia.Z
		delete(procCtx.Assets, termSpec.ValPH)
		return nil
	case procdef.RecvSpec:
		// check via
		gotVia, ok := procCtx.Liabs[termSpec.CommPH]
		if !ok {
			err := typedef.ErrMissingInCtx(termSpec.CommPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typedef.LolliRec)
		if !ok {
			err := typedef.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check value
		gotVal, ok := procCtx.Assets[termSpec.BindPH]
		if !ok {
			err := procdef.ErrMissingInCtx(termSpec.BindPH)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(gotVal, wantVia.Y)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// check cont
		procCtx.Liabs[termSpec.CommPH] = wantVia.Z
		procCtx.Assets[termSpec.BindPH] = wantVia.Y
		return s.checkState(poolID, procEnv, procCtx, procCfg, termSpec.ContTS)
	case procdef.LabSpec:
		// check via
		gotVia, ok := procCtx.Liabs[termSpec.CommPH]
		if !ok {
			err := typedef.ErrMissingInCtx(termSpec.CommPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typedef.PlusRec)
		if !ok {
			err := typedef.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check label
		choice, ok := wantVia.Zs[termSpec.Label]
		if !ok {
			err := fmt.Errorf("label mismatch: want %v, got %q", maps.Keys(wantVia.Zs), termSpec.Label)
			s.log.Error("checking failed")
			return err
		}
		// no cont to check
		procCtx.Liabs[termSpec.CommPH] = choice
		return nil
	case procdef.CaseSpec:
		// check via
		gotVia, ok := procCtx.Liabs[termSpec.CommPH]
		if !ok {
			err := typedef.ErrMissingInCtx(termSpec.CommPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typedef.WithRec)
		if !ok {
			err := typedef.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check conts
		if len(termSpec.ContTSs) != len(wantVia.Zs) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantVia.Zs), len(termSpec.ContTSs))
			s.log.Error("checking failed")
			return err
		}
		for label, choice := range wantVia.Zs {
			cont, ok := termSpec.ContTSs[label]
			if !ok {
				err := fmt.Errorf("label mismatch: want %q, got nothing", label)
				s.log.Error("checking failed")
				return err
			}
			procCtx.Liabs[termSpec.CommPH] = choice
			err := s.checkState(poolID, procEnv, procCtx, procCfg, cont)
			if err != nil {
				s.log.Error("checking failed")
				return err
			}
		}
		return nil
	case procdef.FwdSpec:
		if len(procCtx.Assets) != 1 {
			err := fmt.Errorf("context mismatch: want 1 item, got %v items", len(procCtx.Assets))
			s.log.Error("checking failed")
			return err
		}
		viaSt, ok := procCtx.Liabs[termSpec.X]
		if !ok {
			err := typedef.ErrMissingInCtx(termSpec.X)
			s.log.Error("checking failed")
			return err
		}
		fwdSt, ok := procCtx.Assets[termSpec.Y]
		if !ok {
			err := procdef.ErrMissingInCtx(termSpec.Y)
			s.log.Error("checking failed")
			return err
		}
		if fwdSt.Pol() != viaSt.Pol() {
			err := typedef.ErrPolarityMismatch(fwdSt, viaSt)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(fwdSt, viaSt)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		delete(procCtx.Liabs, termSpec.X)
		delete(procCtx.Assets, termSpec.Y)
		return nil
	default:
		panic(procdef.ErrTermTypeUnexpected(ts))
	}
}

func (s *service) checkClient(
	poolID identity.ADT,
	procEnv procxec.Env,
	procCtx typedef.Context,
	procCfg procxec.Cfg,
	ts procdef.TermSpec,
) error {
	switch termSpec := ts.(type) {
	case procdef.CloseSpec:
		err := procdef.ErrTermTypeMismatch(ts, procdef.WaitSpec{})
		s.log.Error("checking failed")
		return err
	case procdef.WaitSpec:
		// check via
		gotVia, ok := procCtx.Assets[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCtx(termSpec.CommPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typedef.OneRec)
		if !ok {
			err := typedef.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check cont
		delete(procCtx.Assets, termSpec.CommPH)
		return s.checkState(poolID, procEnv, procCtx, procCfg, termSpec.ContTS)
	case procdef.SendSpec:
		// check via
		gotVia, ok := procCtx.Assets[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCtx(termSpec.CommPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typedef.LolliRec)
		if !ok {
			err := typedef.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check value
		gotVal, ok := procCtx.Assets[termSpec.ValPH]
		if !ok {
			err := procdef.ErrMissingInCtx(termSpec.ValPH)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(gotVal, wantVia.Y)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		procCtx.Assets[termSpec.CommPH] = wantVia.Z
		delete(procCtx.Assets, termSpec.ValPH)
		return nil
	case procdef.RecvSpec:
		// check via
		gotVia, ok := procCtx.Assets[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCtx(termSpec.CommPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typedef.TensorRec)
		if !ok {
			err := typedef.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check value
		gotVal, ok := procCtx.Assets[termSpec.BindPH]
		if !ok {
			err := procdef.ErrMissingInCtx(termSpec.BindPH)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(gotVal, wantVia.Y)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// check cont
		procCtx.Assets[termSpec.CommPH] = wantVia.Z
		procCtx.Assets[termSpec.BindPH] = wantVia.Y
		return s.checkState(poolID, procEnv, procCtx, procCfg, termSpec.ContTS)
	case procdef.LabSpec:
		// check via
		gotVia, ok := procCtx.Assets[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCtx(termSpec.CommPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typedef.WithRec)
		if !ok {
			err := typedef.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check label
		choice, ok := wantVia.Zs[termSpec.Label]
		if !ok {
			err := fmt.Errorf("label mismatch: want %v, got %q", maps.Keys(wantVia.Zs), termSpec.Label)
			s.log.Error("checking failed")
			return err
		}
		procCtx.Assets[termSpec.CommPH] = choice
		return nil
	case procdef.CaseSpec:
		// check via
		gotVia, ok := procCtx.Assets[termSpec.CommPH]
		if !ok {
			err := procdef.ErrMissingInCtx(termSpec.CommPH)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := gotVia.(typedef.PlusRec)
		if !ok {
			err := typedef.ErrSnapTypeMismatch(gotVia, wantVia)
			s.log.Error("checking failed")
			return err
		}
		// check conts
		if len(termSpec.ContTSs) != len(wantVia.Zs) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantVia.Zs), len(termSpec.ContTSs))
			s.log.Error("checking failed")
			return err
		}
		for label, choice := range wantVia.Zs {
			cont, ok := termSpec.ContTSs[label]
			if !ok {
				err := fmt.Errorf("label mismatch: want %q, got nothing", label)
				s.log.Error("checking failed")
				return err
			}
			procCtx.Assets[termSpec.CommPH] = choice
			err := s.checkState(poolID, procEnv, procCtx, procCfg, cont)
			if err != nil {
				s.log.Error("checking failed")
				return err
			}
		}
		return nil
	case procdef.SpawnSpecOld:
		procSig, ok := procEnv.ProcDecs[termSpec.SigID]
		if !ok {
			err := procdec.ErrRootMissingInEnv(termSpec.SigID)
			s.log.Error("checking failed")
			return err
		}
		// check vals
		if len(termSpec.Ys) != len(procSig.Ys) {
			err := fmt.Errorf("context mismatch: want %v items, got %v items", len(procSig.Ys), len(termSpec.Ys))
			s.log.Error("checking failed", slog.Any("want", procSig.Ys), slog.Any("got", termSpec.Ys))
			return err
		}
		if len(termSpec.Ys) == 0 {
			return nil
		}
		for i, ep := range procSig.Ys {
			valRole, ok := procEnv.TypeDefs[ep.TypeQN]
			if !ok {
				err := typedef.ErrSymMissingInEnv(ep.TypeQN)
				s.log.Error("checking failed")
				return err
			}
			wantVal, ok := procEnv.TypeTerms[valRole.TermID]
			if !ok {
				err := typedef.ErrMissingInEnv(valRole.TermID)
				s.log.Error("checking failed")
				return err
			}
			gotVal, ok := procCtx.Assets[termSpec.Ys[i]]
			if !ok {
				err := procdef.ErrMissingInCtx(ep.BindPH)
				s.log.Error("checking failed")
				return err
			}
			err := typedef.CheckRec(gotVal, wantVal)
			if err != nil {
				s.log.Error("checking failed", slog.Any("want", wantVal), slog.Any("got", gotVal))
				return err
			}
			delete(procCtx.Assets, termSpec.Ys[i])
		}
		// check via
		viaRole, ok := procEnv.TypeDefs[procSig.X.TypeQN]
		if !ok {
			err := typedef.ErrSymMissingInEnv(procSig.X.TypeQN)
			s.log.Error("checking failed")
			return err
		}
		wantVia, ok := procEnv.TypeTerms[viaRole.TermID]
		if !ok {
			err := typedef.ErrMissingInEnv(viaRole.TermID)
			s.log.Error("checking failed")
			return err
		}
		// check cont
		procCtx.Assets[termSpec.X] = wantVia
		return s.checkState(poolID, procEnv, procCtx, procCfg, termSpec.ContTS)
	default:
		panic(procdef.ErrTermTypeUnexpected(ts))
	}
}

func errOptimisticUpdate(got revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}

func errMissingPool(want qualsym.ADT) error {
	return fmt.Errorf("pool missing in env: %v", want)
}

func errMissingSig(want identity.ADT) error {
	return fmt.Errorf("sig missing in env: %v", want)
}

func errMissingRole(want qualsym.ADT) error {
	return fmt.Errorf("role missing in env: %v", want)
}
