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
	"orglang/go-engine/adt/option"
	"orglang/go-engine/adt/polarity"
	"orglang/go-engine/adt/procconn"
	"orglang/go-engine/adt/procdec"
	"orglang/go-engine/adt/procdef"
	"orglang/go-engine/adt/procexp"
	"orglang/go-engine/adt/procstep"
	"orglang/go-engine/adt/seqnum"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/typedef"
	"orglang/go-engine/adt/typeexp"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

type API interface {
	Take(procstep.StepSpec) error
	RetrieveSnap(implsem.SemRef) (ExecSnap, error)
}

type ExecRec struct {
	ImplRef implsem.SemRef
	ChnlPH  symbol.ADT
	ChnlON  seqnum.ADT
}

// aka Configuration
type ExecSnap struct {
	ImplRef    implsem.SemRef
	LinearVars map[symbol.ADT]implvar.LinearRec
}

type Env struct {
	TypeDefs map[uniqsym.ADT]typedef.DefRec
	TypeExps map[valkey.ADT]typeexp.ExpRec
	ProcDecs map[identity.ADT]procdec.DecRec
}

func ChnlPH(rec implvar.LinearRec) symbol.ADT { return rec.ChnlPH }

type ExecMod struct {
	ImplRefs   []implsem.SemRef
	LinearVars []implvar.LinearRec
}

type service struct {
	procExecs Repo
	procConns procconn.Repo
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
	procConns procconn.Repo,
	procDecs procdec.Repo,
	typeDefs typedef.Repo,
	typeExps typeexp.Repo,
	operator db.Operator,
	l *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{procExecs, procConns, procDecs, typeDefs, typeExps, operator, l.With(name)}
}

func (s *service) RetrieveSnap(ref implsem.SemRef) (_ ExecSnap, err error) {
	return ExecSnap{}, nil
}

func ErrMissingChnl(want symbol.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func (s *service) Take(spec procstep.StepSpec) (err error) {
	implAttr := slog.Any("proc", spec.ImplRef)
	s.log.Debug("step taking started", implAttr, slog.Any("exp", spec.ProcExp))
	ctx := context.Background()
	// initial values
	implRef := spec.ImplRef
	expSpec := spec.ProcExp
	for expSpec != nil {
		var execSnap ExecSnap
		selectErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
			execSnap, err = s.procExecs.SelectSnap(ds, implRef)
			return err
		})
		if selectErr1 != nil {
			s.log.Error("step taking failed", implAttr)
			return selectErr1
		}
		if len(execSnap.LinearVars) == 0 {
			panic("zero channel binds")
		}
		decIDs := procexp.CollectEnv(expSpec)
		var procDecs map[identity.ADT]procdec.DecRec
		selectErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
			procDecs, err = s.procDecs.SelectEnv(ds, decIDs)
			return err
		})
		if selectErr2 != nil {
			s.log.Error("step taking failed", implAttr, slog.Any("decs", decIDs))
			return selectErr2
		}
		typeQNs := procdec.CollectEnv(maps.Values(procDecs))
		var typeDefs map[uniqsym.ADT]typedef.DefRec
		selectErr3 := s.operator.Implicit(ctx, func(ds db.Source) error {
			typeDefs, err = s.typeDefs.SelectEnv(ds, typeQNs)
			return err
		})
		if selectErr3 != nil {
			s.log.Error("step taking failed", implAttr, slog.Any("types", typeQNs))
			return selectErr3
		}
		envIDs := typedef.CollectEnv(maps.Values(typeDefs))
		ctxIDs := CollectCtx(maps.Values(execSnap.LinearVars))
		var typeExps map[valkey.ADT]typeexp.ExpRec
		selectErr4 := s.operator.Implicit(ctx, func(ds db.Source) error {
			typeExps, err = s.typeExps.SelectEnv(ds, append(envIDs, ctxIDs...))
			return err
		})
		if selectErr4 != nil {
			s.log.Error("step taking failed", implAttr, slog.Any("env", envIDs), slog.Any("ctx", ctxIDs))
			return selectErr4
		}
		procEnv := Env{ProcDecs: procDecs, TypeDefs: typeDefs, TypeExps: typeExps}
		procCtx := convertToCtx(maps.Values(execSnap.LinearVars), typeExps)
		// type checking
		err = s.checkType(procEnv, procCtx, execSnap, expSpec)
		if err != nil {
			s.log.Error("step taking failed", implAttr)
			return err
		}
		// step taking
		nextSpec, _, procMod, err := s.takeWith(procEnv, execSnap, expSpec)
		if err != nil {
			s.log.Error("step taking failed", implAttr)
			return err
		}
		err = s.operator.Explicit(ctx, func(ds db.Source) error {
			err = s.procExecs.UpdateProc(ds, procMod)
			if err != nil {
				s.log.Error("step taking failed", implAttr)
				return err
			}
			return nil
		})
		if err != nil {
			s.log.Error("step taking failed", implAttr)
			return err
		}
		// next values
		implRef = nextSpec.ImplRef
		expSpec = nextSpec.ProcExp
	}
	s.log.Debug("step taking succeed", implAttr)
	return nil
}

func (s *service) takeWith(
	procEnv Env,
	execSnap ExecSnap,
	es procexp.ExpSpec,
) (
	stepSpec procstep.StepSpec,
	connMod procconn.ConnMod,
	execMod ExecMod,
	err error,
) {
	ctx := context.Background()
	implAttr := slog.Any("impl", execSnap.ImplRef)
	execMod.ImplRefs = append(execMod.ImplRefs, execSnap.ImplRef)
	switch expSpec := es.(type) {
	case procexp.CloseSpec:
		commChnl, ok := execSnap.LinearVars[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implAttr)
			return procstep.StepSpec{}, connMod, ExecMod{}, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		commAttr := slog.Any("comm", commChnl.CommRef)
		// получаем снепшот соединения
		var connSnap procconn.ConnSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.procConns.SelectSnapByQry(ds, procconn.ConnQuery{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implAttr, commAttr)
			return stepSpec, connMod, execMod, getErr
		}
		// обнуляем канал закрывателя
		execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
			ImplRef: commChnl.ImplRef,
			CommRef: commChnl.CommRef,
			ChnlID:  identity.Nil,
			ChnlPH:  commChnl.ChnlPH,
			ChnlBS:  commChnl.ChnlBS,
			ExpVK:   valkey.Zero,
		})
		subscription := connSnap.NextStep()
		if subscription == nil {
			// регистрируем сообщение закрывателя
			connMod.Steps = append(connMod.Steps, procstep.PubRec{
				CommRef: connSnap.CommRef,
				ImplRef: execSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ValExp:  procexp.CloseRec(expSpec),
			})
			s.log.Debug("taking half done", implAttr, commAttr)
			return stepSpec, connMod, execMod, nil
		}
		observation, ok := subscription.(procstep.SubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(subscription))
		}
		switch contExp := observation.ContExp.(type) {
		case procexp.WaitRec:
			stepSpec = procstep.StepSpec{
				ImplRef: observation.ImplRef,
				ProcExp: contExp.ContES,
			}
			s.log.Debug("step taking succeed", implAttr, commAttr)
			return stepSpec, connMod, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(observation.ContExp))
		}
	case procexp.WaitSpec:
		commChnl, ok := execSnap.LinearVars[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implAttr)
			return procstep.StepSpec{}, connMod, ExecMod{}, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		// обнуляем канал наблюдателя
		execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
			ImplRef: commChnl.ImplRef,
			CommRef: commChnl.CommRef,
			ChnlID:  identity.Nil,
			ChnlPH:  commChnl.ChnlPH,
			ChnlBS:  commChnl.ChnlBS,
			ExpVK:   valkey.Zero,
		})
		// получаем снепшот соединения
		var connSnap procconn.ConnSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.procConns.SelectSnapByQry(ds, procconn.ConnQuery{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implAttr)
			return stepSpec, connMod, execMod, getErr
		}
		publication := connSnap.NextStep()
		if publication == nil {
			// регистрируем подписку наблюдателя
			connMod.Steps = append(connMod.Steps, procstep.SubRec{
				CommRef: connSnap.CommRef,
				ImplRef: execSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: procexp.WaitRec(expSpec),
			})
			s.log.Debug("taking half done", implAttr)
			return stepSpec, connMod, execMod, nil
		}
		closage, ok := publication.(procstep.PubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(publication))
		}
		switch valExp := closage.ValExp.(type) {
		case procexp.CloseRec:
			stepSpec = procstep.StepSpec{
				ImplRef: execSnap.ImplRef,
				ProcExp: expSpec.ContES,
			}
			s.log.Debug("step taking succeed", implAttr)
			return stepSpec, connMod, execMod, nil
		case procexp.FwdRec:
			// перенаправляем продолжение наблюдателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: execSnap.ImplRef,
				CommRef: connSnap.CommRef,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   commChnl.ExpVK,
			})
			stepSpec = procstep.StepSpec{
				ImplRef: execSnap.ImplRef,
				ProcExp: expSpec,
			}
			s.log.Debug("step taking succeed", implAttr)
			return stepSpec, connMod, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(closage.ValExp))
		}
	case procexp.SendSpec:
		commChnl, ok := execSnap.LinearVars[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("step taking failed", implAttr)
			return procstep.StepSpec{}, connMod, ExecMod{}, err
		}
		// лишаем значения отправителя
		execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
			ImplRef: commChnl.ImplRef,
			CommRef: commChnl.CommRef,
			ChnlID:  commChnl.ChnlID,
			ChnlPH:  expSpec.ValChnlPH,
			ChnlBS:  commChnl.ChnlBS,
			ExpVK:   commChnl.ExpVK.Invert(),
		})
		typeExp, ok := procEnv.TypeExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", implAttr)
			return procstep.StepSpec{}, connMod, ExecMod{}, typedef.ErrMissingInEnv(commChnl.ExpVK)
		}
		nextExpVK := typeExp.(typeexp.ProdRec).Next()
		valChnl, ok := execSnap.LinearVars[expSpec.ValChnlPH]
		if !ok {
			s.log.Error("step taking failed", implAttr)
			return procstep.StepSpec{}, connMod, ExecMod{}, procdef.ErrMissingInCfg(expSpec.ValChnlPH)
		}
		// получаем снепшот соединения
		var connSnap procconn.ConnSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.procConns.SelectSnapByQry(ds, procconn.ConnQuery{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implAttr)
			return stepSpec, connMod, execMod, getErr
		}
		subscription := connSnap.NextStep()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение отправителя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем сообщение отправителя
			connMod.Steps = append(connMod.Steps, procstep.PubRec{
				CommRef: connSnap.CommRef,
				ImplRef: execSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: procexp.SendRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					ValChnlID:  valChnl.ChnlID,
					ValExpVK:   valChnl.ExpVK,
				},
			})
			s.log.Debug("taking half done", implAttr)
			return stepSpec, connMod, execMod, nil
		}
		receival, ok := subscription.(procstep.SubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(subscription))
		}
		switch contExp := receival.ContExp.(type) {
		case procexp.RecvRec:
			// вяжем продолжение отправителя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// вяжем значение принимателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: receival.ImplRef,
				CommRef: valChnl.CommRef,
				ChnlID:  valChnl.ChnlID,
				ChnlPH:  contExp.NewChnlPH,
				ChnlBS:  valChnl.ChnlBS,
				ExpVK:   valChnl.ExpVK,
			})
			stepSpec = procstep.StepSpec{
				ImplRef: receival.ImplRef,
				ProcExp: contExp.ContES,
			}
			s.log.Debug("step taking succeed", implAttr)
			return stepSpec, connMod, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(receival.ContExp))
		}
	case procexp.RecvSpec:
		commChnl, ok := execSnap.LinearVars[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", implAttr)
			return procstep.StepSpec{}, connMod, ExecMod{}, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		typeExp, ok := procEnv.TypeExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", implAttr)
			return procstep.StepSpec{}, connMod, ExecMod{}, typedef.ErrMissingInEnv(commChnl.ExpVK)
		}
		nextExpVK := typeExp.(typeexp.ProdRec).Next()
		// получаем снепшот соединения
		var connSnap procconn.ConnSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.procConns.SelectSnapByQry(ds, procconn.ConnQuery{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implAttr)
			return stepSpec, connMod, execMod, getErr
		}
		publication := connSnap.NextStep()
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение принимателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем подписку принимателя
			connMod.Steps = append(connMod.Steps, procstep.SubRec{
				CommRef: connSnap.CommRef,
				ImplRef: execSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: procexp.RecvRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					NewChnlPH:  expSpec.NewChnlPH,
					ContES:     expSpec.ContES,
				},
			})
			s.log.Debug("taking half done", implAttr)
			return stepSpec, connMod, execMod, nil
		}
		sending, ok := publication.(procstep.PubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(publication))
		}
		switch valExp := sending.ValExp.(type) {
		case procexp.SendRec:
			// вяжем продолжение принимателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// вяжем значение принимателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: valExp.CommRef,
				ChnlID:  valExp.ValChnlID,
				ChnlPH:  expSpec.NewChnlPH,
				ChnlBS:  implvar.AssetSide,
				ExpVK:   valExp.ValExpVK,
			})
			stepSpec = procstep.StepSpec{
				ImplRef: execSnap.ImplRef,
				ProcExp: expSpec.ContES,
			}
			s.log.Debug("step taking succeed", implAttr)
			return stepSpec, connMod, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(sending.ValExp))
		}
	case procexp.LabSpec:
		commChnl, ok := execSnap.LinearVars[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("step taking failed")
			return procstep.StepSpec{}, connMod, ExecMod{}, err
		}
		typeExp, ok := procEnv.TypeExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", implAttr)
			return procstep.StepSpec{}, connMod, ExecMod{}, typedef.ErrMissingInEnv(commChnl.ExpVK)
		}
		nextExpVK := typeExp.(typeexp.SumRec).Next(expSpec.ValLabQN)
		// получаем снепшот соединения
		var connSnap procconn.ConnSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.procConns.SelectSnapByQry(ds, procconn.ConnQuery{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implAttr)
			return stepSpec, connMod, execMod, getErr
		}
		subscription := connSnap.NextStep()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение решателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// регистрируем сообщение решателя
			connMod.Steps = append(connMod.Steps, procstep.PubRec{
				CommRef: connSnap.CommRef,
				ImplRef: execSnap.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: procexp.LabRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					ValLabQN:   expSpec.ValLabQN,
				},
			})
			s.log.Debug("taking half done", implAttr)
			return stepSpec, connMod, execMod, nil
		}
		folowing, ok := subscription.(procstep.SubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(subscription))
		}
		switch contExp := folowing.ContExp.(type) {
		case procexp.CaseRec:
			// вяжем продолжение решателя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpVK:   nextExpVK,
			})
			// вяжем продолжение последователя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: folowing.ImplRef,
				CommRef: folowing.CommRef,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  contExp.CommChnlPH,
				// TODO значение ChnlBS
				ExpVK: nextExpVK,
			})
			stepSpec = procstep.StepSpec{
				ImplRef: folowing.ImplRef,
				ProcExp: contExp.ContESs[expSpec.ValLabQN],
			}
			s.log.Debug("step taking succeed", implAttr)
			return stepSpec, connMod, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(folowing.ContExp))
		}
	case procexp.CaseSpec:
		commChnl, ok := execSnap.LinearVars[expSpec.CommChnlPH]
		if !ok {
			err := procdef.ErrMissingInCfg(expSpec.CommChnlPH)
			s.log.Error("step taking failed")
			return procstep.StepSpec{}, connMod, ExecMod{}, err
		}
		// получаем снепшот соединения
		var connSnap procconn.ConnSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.procConns.SelectSnapByQry(ds, procconn.ConnQuery{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implAttr)
			return stepSpec, connMod, execMod, getErr
		}
		publication := connSnap.NextStep()
		if publication == nil {
			newChnlID := identity.New()
			// регистрируем подписку последователя
			connMod.Steps = append(connMod.Steps, procstep.SubRec{
				CommRef: commChnl.CommRef,
				ImplRef: commChnl.ImplRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: procexp.CaseRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					ContESs:    expSpec.ContESes,
				},
			})
			s.log.Debug("taking half done", implAttr)
			return stepSpec, connMod, execMod, nil
		}
		decision, ok := publication.(procstep.PubRec)
		if !ok {
			panic(procstep.ErrRecTypeUnexpected(publication))
		}
		switch valExp := decision.ValExp.(type) {
		case procexp.LabRec:
			// вяжем продолжение последователя
			execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
				ImplRef: commChnl.ImplRef,
				CommRef: commChnl.CommRef,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				// TODO значение ChnlBS
				ExpVK: valExp.ValExpVK,
			})
			stepSpec = procstep.StepSpec{
				ImplRef: execSnap.ImplRef,
				ProcExp: expSpec.ContESes[valExp.ValLabQN],
			}
			s.log.Debug("step taking succeed", implAttr)
			return stepSpec, connMod, execMod, nil
		default:
			panic(procexp.ErrRecTypeUnexpected(decision.ValExp))
		}
	case procexp.FwdSpec:
		commChnl, ok := execSnap.LinearVars[expSpec.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed")
			return procstep.StepSpec{}, connMod, ExecMod{}, procdef.ErrMissingInCfg(expSpec.CommChnlPH)
		}
		contChnl, ok := execSnap.LinearVars[expSpec.ContChnlPH]
		if !ok {
			s.log.Error("step taking failed")
			return procstep.StepSpec{}, connMod, ExecMod{}, procdef.ErrMissingInCfg(expSpec.ContChnlPH)
		}
		typeExp, ok := procEnv.TypeExps[commChnl.ExpVK]
		if !ok {
			s.log.Error("step taking failed", implAttr)
			return procstep.StepSpec{}, connMod, ExecMod{}, typedef.ErrMissingInEnv(commChnl.ExpVK)
		}
		// получаем снепшот соединения
		var connSnap procconn.ConnSnap
		getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
			connSnap, err = s.procConns.SelectSnapByQry(ds, procconn.ConnQuery{
				CommRef: commChnl.CommRef,
				ChnlID:  option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", implAttr)
			return stepSpec, connMod, execMod, getErr
		}
		communication := connSnap.NextStep()
		switch typeExp.Pol() {
		case polarity.Pos:
			switch forwardable := communication.(type) {
			case procstep.SubRec:
				// перенаправляем подписчика
				execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
					ImplRef: forwardable.ImplRef,
					CommRef: forwardable.CommRef,
					ChnlID:  commChnl.ChnlID,
					ChnlPH:  forwardable.ContExp.Via(),
					// TODO значение ChnlBS
					ExpVK: commChnl.ExpVK,
				})
				stepSpec = procstep.StepSpec{
					ImplRef: forwardable.ImplRef,
					ProcExp: forwardable.ContExp,
				}
				s.log.Debug("step taking succeed", implAttr)
				return stepSpec, connMod, execMod, nil
			case procstep.PubRec:
				// перенаправляем публикатора
				execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
					ImplRef: forwardable.ImplRef,
					CommRef: forwardable.CommRef,
					ChnlID:  contChnl.ChnlID,
					ChnlPH:  forwardable.ValExp.Via(),
					// TODO значение ChnlBS
					ExpVK: contChnl.ExpVK,
				})
				stepSpec = procstep.StepSpec{
					ImplRef: forwardable.ImplRef,
					ProcExp: forwardable.ValExp,
				}
				s.log.Debug("step taking succeed", implAttr)
				return stepSpec, connMod, execMod, nil
			case nil:
				// лишаем значений схлопывающегося
				execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
					ImplRef: commChnl.ImplRef,
					CommRef: commChnl.CommRef,
					ChnlID:  commChnl.ChnlID,
					ChnlPH:  expSpec.CommChnlPH,
					ChnlBS:  commChnl.ChnlBS,
					ExpVK:   commChnl.ExpVK.Invert(),
				})
				execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
					ImplRef: commChnl.ImplRef,
					// TODO значение CommRef
					ChnlID: contChnl.ChnlID,
					ChnlPH: expSpec.ContChnlPH,
					ChnlBS: contChnl.ChnlBS,
					ExpVK:  contChnl.ExpVK.Invert(),
				})
				// регистрируем сообщение схлопывающегося
				connMod.Steps = append(connMod.Steps, procstep.PubRec{
					CommRef: connSnap.CommRef,
					ImplRef: execSnap.ImplRef,
					ChnlID:  commChnl.ChnlID,
					ValExp: procexp.FwdRec{
						ContChnlID: contChnl.ChnlID,
					},
				})
				s.log.Debug("taking half done", implAttr)
				return stepSpec, connMod, execMod, nil
			default:
				panic(procstep.ErrRecTypeUnexpected(communication))
			}
		case polarity.Neg:
			switch forwardable := communication.(type) {
			case procstep.SubRec:
				// перенаправляем подписчика
				execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
					ImplRef: forwardable.ImplRef,
					CommRef: forwardable.CommRef,
					ChnlID:  contChnl.ChnlID,
					ChnlPH:  forwardable.ContExp.Via(),
					// TODO значение ChnlBS
					ExpVK: contChnl.ExpVK,
				})
				stepSpec = procstep.StepSpec{
					ImplRef: forwardable.ImplRef,
					ProcExp: forwardable.ContExp,
				}
				s.log.Debug("step taking succeed", implAttr)
				return stepSpec, connMod, execMod, nil
			case procstep.PubRec:
				// перенаправляем публикатора
				execMod.LinearVars = append(execMod.LinearVars, implvar.LinearRec{
					ImplRef: forwardable.ImplRef,
					CommRef: forwardable.CommRef,
					ChnlID:  commChnl.ChnlID,
					ChnlPH:  forwardable.ValExp.Via(),
					// TODO значение ChnlBS
					ExpVK: commChnl.ExpVK,
				})
				stepSpec = procstep.StepSpec{
					ImplRef: forwardable.ImplRef,
					ProcExp: forwardable.ValExp, // TODO: несовпадение типов
				}
				s.log.Debug("step taking succeed", implAttr)
				return stepSpec, connMod, execMod, nil
			case nil:
				// регистрируем подписку схлопывающегося
				connMod.Steps = append(connMod.Steps, procstep.SubRec{
					CommRef: connSnap.CommRef,
					ImplRef: commChnl.ImplRef,
					ChnlID:  commChnl.ChnlID,
					ContExp: procexp.FwdRec{
						ContChnlID: contChnl.ChnlID,
					},
				})
				s.log.Debug("taking half done", implAttr)
				return stepSpec, connMod, execMod, nil
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
		if bind.ChnlBS == implvar.LiabSide {
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
	chnlBR, ok := execSnap.LinearVars[expSpec.Via()]
	if !ok {
		panic("no comm chnl in proc snap")
	}
	if chnlBR.ChnlBS == implvar.LiabSide {
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
		if len(expSpec.ContESes) != len(wantVia.Choices) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantVia.Choices), len(expSpec.ContESes))
			s.log.Error("checking failed")
			return err
		}
		for label, choice := range wantVia.Choices {
			cont, ok := expSpec.ContESes[label]
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
		if len(expSpec.ContESes) != len(wantVia.Choices) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantVia.Choices), len(expSpec.ContESes))
			s.log.Error("checking failed")
			return err
		}
		for label, choice := range wantVia.Choices {
			cont, ok := expSpec.ContESes[label]
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

func errOptimisticUpdate(got seqnum.ADT) error {
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
