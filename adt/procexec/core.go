package procexec

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"maps"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/polarity"
	"orglang/go-engine/adt/procdec"
	"orglang/go-engine/adt/procdef"
	"orglang/go-engine/adt/procexp"
	"orglang/go-engine/adt/procstep"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/typedef"
	"orglang/go-engine/adt/typeexp"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

type API interface {
	Take(procstep.CommSpec) error
	RetrieveSnap(implsem.SemRef) (ExecSnap, error)
}

type ExecRec struct {
	ImplRef implsem.SemRef
	ChnlPH  symbol.ADT
}

// aka Configuration
type ExecSnap struct {
	ImplRef implsem.SemRef
	ChnlVRs map[symbol.ADT]implvar.LinearRec
	ProcSRs map[identity.ADT]procstep.CommRec
}

type Env struct {
	TypeDefs map[uniqsym.ADT]typedef.DefRec
	TypeExps map[valkey.ADT]typeexp.ExpRec
	ProcDecs map[identity.ADT]procdec.DecRec
}

func ChnlPH(rec implvar.LinearRec) symbol.ADT { return rec.ChnlPH }

type CommMod struct {
	Refs  []implsem.SemRef
	Vars  []implvar.LinearRec
	Comms []procstep.CommRec
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
	return new(service)
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

func (s *service) RetrieveSnap(ref implsem.SemRef) (_ ExecSnap, err error) {
	return ExecSnap{}, nil
}

func ErrMissingChnl(want symbol.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func (s *service) Take(spec procstep.CommSpec) (err error) {
	procAttr := slog.Any("proc", spec.ExecRef)
	s.log.Debug("step taking started", procAttr, slog.Any("exp", spec.ProcES))
	ctx := context.Background()
	// initial values
	implRef := spec.ExecRef
	expSpec := spec.ProcES
	for expSpec != nil {
		var execSnap ExecSnap
		selectErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
			execSnap, err = s.procExecs.SelectSnap(ds, implRef)
			return err
		})
		if selectErr1 != nil {
			s.log.Error("step taking failed", procAttr)
			return selectErr1
		}
		if len(execSnap.ChnlVRs) == 0 {
			panic("zero channel binds")
		}
		decIDs := procexp.CollectEnv(expSpec)
		var procDecs map[identity.ADT]procdec.DecRec
		selectErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
			procDecs, err = s.procDecs.SelectEnv(ds, decIDs)
			return err
		})
		if selectErr2 != nil {
			s.log.Error("step taking failed", procAttr, slog.Any("decs", decIDs))
			return selectErr2
		}
		typeQNs := procdec.CollectEnv(maps.Values(procDecs))
		var typeDefs map[uniqsym.ADT]typedef.DefRec
		selectErr3 := s.operator.Implicit(ctx, func(ds db.Source) error {
			typeDefs, err = s.typeDefs.SelectEnv(ds, typeQNs)
			return err
		})
		if selectErr3 != nil {
			s.log.Error("step taking failed", procAttr, slog.Any("types", typeQNs))
			return selectErr3
		}
		envIDs := typedef.CollectEnv(maps.Values(typeDefs))
		ctxIDs := CollectCtx(maps.Values(execSnap.ChnlVRs))
		var typeExps map[valkey.ADT]typeexp.ExpRec
		selectErr4 := s.operator.Implicit(ctx, func(ds db.Source) error {
			typeExps, err = s.typeExps.SelectEnv(ds, append(envIDs, ctxIDs...))
			return err
		})
		if selectErr4 != nil {
			s.log.Error("step taking failed", procAttr, slog.Any("env", envIDs), slog.Any("ctx", ctxIDs))
			return selectErr4
		}
		procEnv := Env{ProcDecs: procDecs, TypeDefs: typeDefs, TypeExps: typeExps}
		procCtx := convertToCtx(maps.Values(execSnap.ChnlVRs), typeExps)
		// type checking
		err = s.checkType(procEnv, procCtx, execSnap, expSpec)
		if err != nil {
			s.log.Error("step taking failed", procAttr)
			return err
		}
		// step taking
		nextSpec, procMod, err := s.takeWith(procEnv, execSnap, expSpec)
		if err != nil {
			s.log.Error("step taking failed", procAttr)
			return err
		}
		err = s.operator.Explicit(ctx, func(ds db.Source) error {
			err = s.procExecs.UpdateProc(ds, procMod)
			if err != nil {
				s.log.Error("step taking failed", procAttr)
				return err
			}
			return nil
		})
		if err != nil {
			s.log.Error("step taking failed", procAttr)
			return err
		}
		// next values
		implRef = nextSpec.ExecRef
		expSpec = nextSpec.ProcES
	}
	s.log.Debug("step taking succeed", procAttr)
	return nil
}

func (s *service) takeWith(
	procEnv Env,
	execSnap ExecSnap,
	es procexp.ExpSpec,
) (
	stepSpec procstep.CommSpec,
	commMod CommMod,
	_ error,
) {
	switch expSpec := es.(type) {
	case procexp.CloseSpec:
		commChnl, ok := execSnap.ChnlVRs[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", slog.Any("proc", execSnap.ImplRef))
			return procstep.CommSpec{}, CommMod{}, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		refAttr := slog.Any("chnl", commChnl.ChnlID)
		commMod.Refs = append(commMod.Refs, execSnap.ImplRef)
		// обнуляем канал закрывателя
		commMod.Vars = append(commMod.Vars, implvar.LinearRec{
			ImplRef: execSnap.ImplRef,
			ChnlID:  commChnl.ChnlID,
			ChnlPH:  commChnl.ChnlPH,
			ExpVK:   valkey.Zero,
		})
		subscription := execSnap.ProcSRs[commChnl.ChnlID]
		if subscription == nil {
			// регистрируем сообщение закрывателя
			commMod.Comms = append(commMod.Comms, procstep.PubRec{
				ImplRef: implsem.SemRef{
					ImplID: execSnap.ImplRef.ImplID,
					ImplRN: execSnap.ImplRef.ImplRN.Next(),
				},
				ChnlID: commChnl.ChnlID,
				ValExp: procexp.CloseRec(expSpec),
			})
			s.log.Debug("taking half done", refAttr)
			return stepSpec, commMod, nil
		}
		waiting, ok := subscription.(procstep.SubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(subscription))
		}
		switch contExp := waiting.ContExp.(type) {
		case procexp.WaitRec:
			stepSpec = procstep.CommSpec{
				ExecRef: waiting.ImplRef,
				ProcES:  contExp.ContES,
			}
			s.log.Debug("step taking succeed", refAttr)
			return stepSpec, commMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(waiting.ContExp))
		}
	case procexp.WaitSpec:
		commChnl, ok := execSnap.ChnlVRs[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", slog.Any("proc", execSnap.ImplRef))
			return procstep.CommSpec{}, CommMod{}, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		refAttr := slog.Any("chnl", commChnl.ChnlID)
		commMod.Refs = append(commMod.Refs, execSnap.ImplRef)
		// обнуляем канал наблюдателя
		commMod.Vars = append(commMod.Vars, implvar.LinearRec{
			ImplRef: execSnap.ImplRef,
			ChnlID:  commChnl.ChnlID,
			ChnlPH:  commChnl.ChnlPH,
			ExpVK:   valkey.Zero,
		})
		publication := execSnap.ProcSRs[commChnl.ChnlID]
		if publication == nil {
			// регистрируем подписку наблюдателя
			commMod.Comms = append(commMod.Comms, procstep.SubRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: procexp.WaitRec(expSpec),
			})
			s.log.Debug("taking half done", refAttr)
			return stepSpec, commMod, nil
		}
		closing, ok := publication.(procstep.PubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(publication))
		}
		switch valExp := closing.ValExp.(type) {
		case procexp.CloseRec:
			stepSpec = procstep.CommSpec{
				ExecRef: execSnap.ImplRef,
				ProcES:  expSpec.ContES,
			}
			s.log.Debug("step taking succeed", refAttr)
			return stepSpec, commMod, nil
		case procexp.FwdRec:
			// перенаправляем продолжение наблюдателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ExpVK:   commChnl.ExpVK,
			})
			stepSpec = procstep.CommSpec{
				ExecRef: execSnap.ImplRef,
				ProcES:  expSpec,
			}
			s.log.Debug("step taking succeed", refAttr)
			return stepSpec, commMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(closing.ValExp))
		}
	case procexp.SendSpec:
		commChnl, ok := execSnap.ChnlVRs[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("step taking failed")
			return procstep.CommSpec{}, CommMod{}, err
		}
		refAttr := slog.Any("ref", commChnl.ChnlID)
		commMod.Refs = append(commMod.Refs, execSnap.ImplRef)
		// лишаем значения отправителя
		commMod.Vars = append(commMod.Vars, implvar.LinearRec{
			ImplRef: execSnap.ImplRef,
			ChnlID:  commChnl.ChnlID,
			ChnlPH:  expSpec.ValChnlPH,
			ChnlBS:  commChnl.ChnlBS.Negate(),
		})
		typeExp, ok := procEnv.TypeExps[commChnl.ExpVK]
		if !ok {
			err := typedef.ErrMissingInEnv(commChnl.ExpVK)
			s.log.Error("step taking failed", refAttr)
			return procstep.CommSpec{}, CommMod{}, err
		}
		nextExpVK := typeExp.(typeexp.ProdRec).Next()
		valueVar, ok := execSnap.ChnlVRs[expSpec.ValChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.ValChnlPH)
			s.log.Error("step taking failed", refAttr)
			return procstep.CommSpec{}, CommMod{}, err
		}
		subscription := execSnap.ProcSRs[commChnl.ChnlID]
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение отправителя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				// TODO указать ChnlBS
				ExpVK: nextExpVK,
			})
			// регистрируем сообщение отправителя
			commMod.Comms = append(commMod.Comms, procstep.PubRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: procexp.SendRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					ValChnlID:  valueVar.ChnlID,
					ValExpVK:   valueVar.ExpVK,
				},
			})
			s.log.Debug("taking half done", refAttr)
			return stepSpec, commMod, nil
		}
		receival, ok := subscription.(procstep.SubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(subscription))
		}
		switch contExp := receival.ContExp.(type) {
		case procexp.RecvRec:
			// вяжем продолжение отправителя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// вяжем значение принимателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: receival.ImplRef,
				ChnlID:  valueVar.ChnlID,
				ChnlPH:  contExp.NewChnlPH,
				ChnlBS:  valueVar.ChnlBS,
				ExpVK:   valueVar.ExpVK,
			})
			stepSpec = procstep.CommSpec{
				ExecRef: receival.ImplRef,
				ProcES:  contExp.ContES,
			}
			s.log.Debug("step taking succeed", refAttr)
			return stepSpec, commMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(receival.ContExp))
		}
	case procexp.RecvSpec:
		commChnl, ok := execSnap.ChnlVRs[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("step taking failed")
			return procstep.CommSpec{}, CommMod{}, err
		}
		refAttr := slog.Any("ref", commChnl.ChnlID)
		commMod.Refs = append(commMod.Refs, execSnap.ImplRef)
		typeExp, ok := procEnv.TypeExps[commChnl.ExpVK]
		if !ok {
			err := typedef.ErrMissingInEnv(commChnl.ExpVK)
			s.log.Error("step taking failed", refAttr)
			return procstep.CommSpec{}, CommMod{}, err
		}
		nextExpVK := typeExp.(typeexp.ProdRec).Next()
		publication := execSnap.ProcSRs[commChnl.ChnlID]
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение принимателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем подписку принимателя
			commMod.Comms = append(commMod.Comms, procstep.SubRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: procexp.RecvRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					NewChnlPH:  expSpec.NewChnlPH,
					ContES:     expSpec.ContES,
				},
			})
			s.log.Debug("taking half done", refAttr)
			return stepSpec, commMod, nil
		}
		sending, ok := publication.(procstep.PubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(publication))
		}
		switch valExp := sending.ValExp.(type) {
		case procexp.SendRec:
			// вяжем продолжение принимателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  expSpec.CommChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// вяжем значение принимателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  valExp.ValChnlID,
				ChnlPH:  expSpec.NewChnlPH,
				// TODO значение ChnlBS
				ExpVK: valExp.ValExpVK,
			})
			stepSpec = procstep.CommSpec{
				ExecRef: execSnap.ImplRef,
				ProcES:  expSpec.ContES,
			}
			s.log.Debug("step taking succeed", refAttr)
			return stepSpec, commMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(sending.ValExp))
		}
	case procexp.LabSpec:
		commChnl, ok := execSnap.ChnlVRs[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("step taking failed")
			return procstep.CommSpec{}, CommMod{}, err
		}
		refAttr := slog.Any("ref", commChnl.ChnlID)
		commMod.Refs = append(commMod.Refs, execSnap.ImplRef)
		typeExp, ok := procEnv.TypeExps[commChnl.ExpVK]
		if !ok {
			err := typedef.ErrMissingInEnv(commChnl.ExpVK)
			s.log.Error("step taking failed", refAttr)
			return procstep.CommSpec{}, CommMod{}, err
		}
		nextExpVK := typeExp.(typeexp.SumRec).Next(expSpec.ValLabQN)
		subscription := execSnap.ProcSRs[commChnl.ChnlID]
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение решателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем сообщение решателя
			commMod.Comms = append(commMod.Comms, procstep.PubRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: procexp.LabRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					ValLabQN:   expSpec.ValLabQN,
				},
			})
			s.log.Debug("taking half done", refAttr)
			return stepSpec, commMod, nil
		}
		folowing, ok := subscription.(procstep.SubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(subscription))
		}
		switch contExp := folowing.ContExp.(type) {
		case procexp.CaseRec:
			// вяжем продолжение решателя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// вяжем продолжение последователя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: folowing.ImplRef,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  contExp.CommChnlPH,
				// TODO значение ChnlBS
				ExpVK: nextExpVK,
			})
			stepSpec = procstep.CommSpec{
				ExecRef: folowing.ImplRef,
				ProcES:  contExp.ContESs[expSpec.ValLabQN],
			}
			s.log.Debug("step taking succeed", refAttr)
			return stepSpec, commMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(folowing.ContExp))
		}
	case procexp.CaseSpec:
		commChnl, ok := execSnap.ChnlVRs[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("step taking failed")
			return procstep.CommSpec{}, CommMod{}, err
		}
		refAttr := slog.Any("chnl", commChnl.ChnlID)
		commMod.Refs = append(commMod.Refs, execSnap.ImplRef)
		publication := execSnap.ProcSRs[commChnl.ChnlID]
		if publication == nil {
			newChnlID := identity.New()
			// регистрируем подписку последователя
			commMod.Comms = append(commMod.Comms, procstep.SubRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: procexp.CaseRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					ContESs:    expSpec.ContESs,
				},
			})
			s.log.Debug("taking half done", refAttr)
			return stepSpec, commMod, nil
		}
		decision, ok := publication.(procstep.PubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(publication))
		}
		switch valExp := decision.ValExp.(type) {
		case procexp.LabRec:
			typeExp, ok := procEnv.TypeExps[commChnl.ExpVK]
			if !ok {
				err := typedef.ErrMissingInEnv(commChnl.ExpVK)
				s.log.Error("step taking failed", refAttr)
				return procstep.CommSpec{}, CommMod{}, err
			}
			// вяжем продолжение последователя
			commMod.Vars = append(commMod.Vars, implvar.LinearRec{
				ImplRef: execSnap.ImplRef,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				// TODO значение ChnlBS
				ExpVK: typeExp.(typeexp.SumRec).Next(valExp.ValLabQN),
			})
			stepSpec = procstep.CommSpec{
				ExecRef: execSnap.ImplRef,
				ProcES:  expSpec.ContESs[valExp.ValLabQN],
			}
			s.log.Debug("step taking succeed", refAttr)
			return stepSpec, commMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(decision.ValExp))
		}
	case procexp.FwdSpec:
		commChnl, ok := execSnap.ChnlVRs[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed")
			return procstep.CommSpec{}, CommMod{}, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		contChnl, ok := execSnap.ChnlVRs[expSpec.ContChnlPH]
		if !ok {
			s.log.Error("step taking failed")
			return procstep.CommSpec{}, CommMod{}, procdef.ErrMissingInCfg(expSpec.ContChnlPH)
		}
		refAttr := slog.Any("chnl", commChnl.ChnlID)
		typeExp, ok := procEnv.TypeExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", refAttr)
			return procstep.CommSpec{}, CommMod{}, typedef.ErrMissingInEnv(commChnl.ExpVK)
		}
		communication := execSnap.ProcSRs[commChnl.ChnlID]
		switch typeExp.Pol() {
		case polarity.Pos:
			switch forwardable := communication.(type) {
			case procstep.SubRec:
				// перенаправляем подписчика
				commMod.Vars = append(commMod.Vars, implvar.LinearRec{
					ImplRef: forwardable.ImplRef,
					ChnlID:  commChnl.ChnlID,
					ChnlPH:  forwardable.ContExp.Via(),
					ExpVK:   commChnl.ExpVK,
				})
				stepSpec = procstep.CommSpec{
					ExecRef: forwardable.ImplRef,
					ProcES:  forwardable.ContExp,
				}
				s.log.Debug("step taking succeed", refAttr)
				return stepSpec, commMod, nil
			case procstep.PubRec:
				// перенаправляем публикатора
				commMod.Vars = append(commMod.Vars, implvar.LinearRec{
					ImplRef: forwardable.ImplRef,
					ChnlID:  contChnl.ChnlID,
					ChnlPH:  forwardable.ValExp.Via(),
					ExpVK:   contChnl.ExpVK,
				})
				stepSpec = procstep.CommSpec{
					ExecRef: forwardable.ImplRef,
					ProcES:  forwardable.ValExp,
				}
				s.log.Debug("step taking succeed", refAttr)
				return stepSpec, commMod, nil
			case nil:
				// лишаем значений схлопывающегося
				commMod.Vars = append(commMod.Vars, implvar.LinearRec{
					ImplRef: execSnap.ImplRef,
					ChnlID:  commChnl.ChnlID,
					ChnlBS:  commChnl.ChnlBS.Negate(),
					ChnlPH:  expSpec.CommChnlPH,
				})
				commMod.Vars = append(commMod.Vars, implvar.LinearRec{
					ImplRef: execSnap.ImplRef,
					ChnlID:  contChnl.ChnlID,
					ChnlBS:  contChnl.ChnlBS.Negate(),
					ChnlPH:  expSpec.ContChnlPH,
				})
				// регистрируем сообщение схлопывающегося
				commMod.Comms = append(commMod.Comms, procstep.PubRec{
					ImplRef: execSnap.ImplRef,
					ChnlID:  commChnl.ChnlID,
					ValExp: procexp.FwdRec{
						ContChnlID: contChnl.ChnlID,
					},
				})
				s.log.Debug("taking half done", refAttr)
				return stepSpec, commMod, nil
			default:
				panic(procstep.ErrRecTypeUnexpected(communication))
			}
		case polarity.Neg:
			switch forwardable := communication.(type) {
			case procstep.SubRec:
				// перенаправляем подписчика
				commMod.Vars = append(commMod.Vars, implvar.LinearRec{
					ImplRef: forwardable.ImplRef,
					ChnlID:  contChnl.ChnlID,
					ChnlPH:  forwardable.ContExp.Via(),
					ExpVK:   contChnl.ExpVK,
				})
				stepSpec = procstep.CommSpec{
					ExecRef: forwardable.ImplRef,
					ProcES:  forwardable.ContExp,
				}
				s.log.Debug("step taking succeed", refAttr)
				return stepSpec, commMod, nil
			case procstep.PubRec:
				// перенаправляем публикатора
				commMod.Vars = append(commMod.Vars, implvar.LinearRec{
					ImplRef: forwardable.ImplRef,
					ChnlID:  commChnl.ChnlID,
					ChnlPH:  forwardable.ValExp.Via(),
					ExpVK:   commChnl.ExpVK,
				})
				stepSpec = procstep.CommSpec{
					ExecRef: forwardable.ImplRef,
					ProcES:  forwardable.ValExp, // TODO: несовпадение типов
				}
				s.log.Debug("step taking succeed", refAttr)
				return stepSpec, commMod, nil
			case nil:
				// регистрируем подписку схлопывающегося
				commMod.Comms = append(commMod.Comms, procstep.SubRec{
					ImplRef: execSnap.ImplRef,
					ChnlID:  commChnl.ChnlID,
					ContExp: procexp.FwdRec{
						ContChnlID: contChnl.ChnlID,
					},
				})
				s.log.Debug("taking half done", refAttr)
				return stepSpec, commMod, nil
			default:
				panic(procstep.ErrRecTypeUnexpected(communication))
			}
		default:
			panic(typeexp.ErrPolarityUnexpected(typeExp))
		}
	default:
		panic(procexp.ErrExpTypeUnexpected(es))
	}
}

func CollectCtx(chnls iter.Seq[implvar.LinearRec]) []valkey.ADT {
	return nil
}

func convertToCtx(chnlBinds iter.Seq[implvar.LinearRec], typeExps map[valkey.ADT]typeexp.ExpRec) typedef.Context {
	assets := make(map[symbol.ADT]typeexp.ExpRec, 1)
	liabs := make(map[symbol.ADT]typeexp.ExpRec, 1)
	for bind := range chnlBinds {
		if bind.ChnlBS == implvar.Provider {
			liabs[bind.ChnlPH] = typeExps[bind.ExpVK]
		} else {
			assets[bind.ChnlPH] = typeExps[bind.ExpVK]
		}
	}
	return typedef.Context{Assets: assets, Liabs: liabs}
}

func (s *service) checkType(
	procEnv Env,
	procCtx typedef.Context,
	execSnap ExecSnap,
	expSpec procexp.ExpSpec,
) error {
	chnlBR, ok := execSnap.ChnlVRs[expSpec.Via()]
	if !ok {
		panic("no comm chnl in proc snap")
	}
	if chnlBR.ChnlBS == implvar.Provider {
		return s.checkProvider(procEnv, procCtx, execSnap, expSpec)
	}
	return s.checkClient(procEnv, procCtx, execSnap, expSpec)
}

func (s *service) checkProvider(
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
		err := typeexp.CheckRec(gotVal, wantVia.Val)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// no cont to check
		procCtx.Liabs[expSpec.CommChnlPH] = wantVia.Cont
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
		gotVal, ok := procCtx.Assets[expSpec.NewChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.NewChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typeexp.CheckRec(gotVal, wantVia.Val)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// check cont
		procCtx.Liabs[expSpec.CommChnlPH] = wantVia.Cont
		procCtx.Assets[expSpec.NewChnlPH] = wantVia.Val
		return s.checkType(procEnv, procCtx, procCfg, expSpec.ContES)
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
		choice, ok := wantVia.Choices[expSpec.ValLabQN]
		if !ok {
			err := fmt.Errorf("label mismatch: want %v, got %v", maps.Keys(wantVia.Choices), expSpec.ValLabQN)
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
		if len(expSpec.ContESs) != len(wantVia.Choices) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantVia.Choices), len(expSpec.ContESs))
			s.log.Error("checking failed")
			return err
		}
		for label, choice := range wantVia.Choices {
			cont, ok := expSpec.ContESs[label]
			if !ok {
				err := fmt.Errorf("label mismatch: want %v, got nothing", label)
				s.log.Error("checking failed")
				return err
			}
			procCtx.Liabs[expSpec.CommChnlPH] = choice
			err := s.checkType(procEnv, procCtx, procCfg, cont)
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
		return s.checkType(procEnv, procCtx, procCfg, expSpec.ContES)
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
		err := typeexp.CheckRec(gotVal, wantVia.Val)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		procCtx.Assets[expSpec.CommChnlPH] = wantVia.Cont
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
		gotVal, ok := procCtx.Assets[expSpec.NewChnlPH]
		if !ok {
			err := procdef.ErrMissingInCtx(expSpec.NewChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typeexp.CheckRec(gotVal, wantVia.Val)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// check cont
		procCtx.Assets[expSpec.CommChnlPH] = wantVia.Cont
		procCtx.Assets[expSpec.NewChnlPH] = wantVia.Val
		return s.checkType(procEnv, procCtx, procCfg, expSpec.ContES)
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
		choice, ok := wantVia.Choices[expSpec.ValLabQN]
		if !ok {
			err := fmt.Errorf("label mismatch: want %v, got %v", maps.Keys(wantVia.Choices), expSpec.ValLabQN)
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
		if len(expSpec.ContESs) != len(wantVia.Choices) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantVia.Choices), len(expSpec.ContESs))
			s.log.Error("checking failed")
			return err
		}
		for label, choice := range wantVia.Choices {
			cont, ok := expSpec.ContESs[label]
			if !ok {
				err := fmt.Errorf("label mismatch: want %v, got nothing", label)
				s.log.Error("checking failed")
				return err
			}
			procCtx.Assets[expSpec.CommChnlPH] = choice
			err := s.checkType(procEnv, procCtx, procCfg, cont)
			if err != nil {
				s.log.Error("checking failed")
				return err
			}
		}
		return nil
	// case procexp.SpawnSpecOld:
	// 	procDec, ok := procEnv.ProcDecs[expSpec.SigID]
	// 	if !ok {
	// 		err := procdec.ErrRootMissingInEnv(expSpec.SigID)
	// 		s.log.Error("checking failed")
	// 		return err
	// 	}
	// 	// check vals
	// 	if len(expSpec.Ys) != len(procDec.ClientVRs) {
	// 		err := fmt.Errorf("context mismatch: want %v items, got %v items", len(procDec.ClientVRs), len(expSpec.Ys))
	// 		s.log.Error("checking failed", slog.Any("want", procDec.ClientVRs), slog.Any("got", expSpec.Ys))
	// 		return err
	// 	}
	// 	if len(expSpec.Ys) == 0 {
	// 		return nil
	// 	}
	// 	for i, ep := range procDec.ClientVRs {
	// 		valType, ok := procEnv.TypeDefs[ep.ImplQN]
	// 		if !ok {
	// 			err := typedef.ErrSymMissingInEnv(ep.ImplQN)
	// 			s.log.Error("checking failed")
	// 			return err
	// 		}
	// 		wantVal, ok := procEnv.TypeExps[valType.ExpVK]
	// 		if !ok {
	// 			err := typedef.ErrMissingInEnv(valType.ExpVK)
	// 			s.log.Error("checking failed")
	// 			return err
	// 		}
	// 		gotVal, ok := procCtx.Assets[expSpec.Ys[i]]
	// 		if !ok {
	// 			err := procdef.ErrMissingInCtx(ep.ChnlPH)
	// 			s.log.Error("checking failed")
	// 			return err
	// 		}
	// 		err := typeexp.CheckRec(gotVal, wantVal)
	// 		if err != nil {
	// 			s.log.Error("checking failed", slog.Any("want", wantVal), slog.Any("got", gotVal))
	// 			return err
	// 		}
	// 		delete(procCtx.Assets, expSpec.Ys[i])
	// 	}
	// 	// check via
	// 	viaRole, ok := procEnv.TypeDefs[procDec.ProviderVR.ImplQN]
	// 	if !ok {
	// 		err := typedef.ErrSymMissingInEnv(procDec.ProviderVR.ImplQN)
	// 		s.log.Error("checking failed")
	// 		return err
	// 	}
	// 	wantVia, ok := procEnv.TypeExps[viaRole.ExpVK]
	// 	if !ok {
	// 		err := typedef.ErrMissingInEnv(viaRole.ExpVK)
	// 		s.log.Error("checking failed")
	// 		return err
	// 	}
	// 	// check cont
	// 	procCtx.Assets[expSpec.X] = wantVia
	// 	return s.checkType(procEnv, procCtx, procCfg, expSpec.ContES)
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
