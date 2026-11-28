package core

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"maps"
	"reflect"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/execvar"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/polarity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	commact "orglang/go-engine/internal/proc/commact/core"
	execcomm "orglang/go-engine/internal/proc/execcomm/core"
	execstep "orglang/go-engine/internal/proc/execstep/core"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	termdef "orglang/go-engine/internal/proc/termdef/core"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type API interface {
	Take(context.Context, execstep.StepSpec) error
	RetrieveSnap(execroot.RootRef) (EvalSnap, error)
}

type EvalRec struct {
	RootRef  execroot.RootRef
	LiabMode execvar.Mode
}

type EvalMod struct {
	ExecRefs   []execroot.RootRef
	LinearVars []execvar.LinearRec
}

type EvalEff struct {
	Steps []execstep.StepSpec
}

type takeResult struct {
	Mod           EvalMod
	Effects       EvalEff
	Communication execcomm.CommMod
}

// aka Configuration
type EvalSnap struct {
	RootRef    execroot.RootRef
	LinearVars map[symbol.ADT]execvar.LinearRec
}

type Env struct {
	TypeDefs map[uniqsym.ADT]typedef.DefRec
	TypeExps map[contkey.ADT]typedef.ExpRec
	ProcDecs map[identity.ADT]termdec.DecRec
}

func ChnlPH(rec execvar.LinearRec) symbol.ADT { return rec.ChnlPH }

type service struct {
	execEvalRepo Repo
	commRepo     execcomm.Repo
	termDecRepo  termdec.Repo
	typeDefRepo  typedef.Repo
	typeExpRepo  typedef.TypeExpRepo
	transactor   tx.Transactor
	log          *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func NewService(
	execEvalRepo Repo,
	commRepo execcomm.Repo,
	termDecRepo termdec.Repo,
	typeDefRepo typedef.Repo,
	typeExpRepo typedef.TypeExpRepo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{
		execEvalRepo, commRepo,
		termDecRepo, typeDefRepo, typeExpRepo,
		transactor, log.With(name),
	}
}

func (s *service) RetrieveSnap(ref execroot.RootRef) (_ EvalSnap, err error) {
	return EvalSnap{}, nil
}

func ErrMissingChnl(want symbol.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func (s *service) Take(ctx context.Context, spec execstep.StepSpec) (err error) {
	execAttr := slog.Any("proc", spec.RootRef)
	s.log.Debug("step taking started", execAttr, slog.Any("exp", spec.ProcExp))
	// initial values
	execRef := spec.RootRef
	expSpec := spec.ProcExp
	for expSpec != nil {
		evalSnap, getErr1 := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (EvalSnap, error) {
			return s.execEvalRepo.GetSnapByRef(txCtx, conn, execRef)
		})
		if getErr1 != nil {
			s.log.Error("step taking failed", execAttr)
			return getErr1
		}
		if len(evalSnap.LinearVars) == 0 {
			panic("zero channel binds")
		}
		decIDs := commact.CollectEnv(expSpec)
		procDecs, getErr2 := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (map[identity.ADT]termdec.DecRec, error) {
			return s.termDecRepo.SelectEnv(txCtx, conn, decIDs)
		})
		if getErr2 != nil {
			s.log.Error("step taking failed", execAttr, slog.Any("decs", decIDs))
			return getErr2
		}
		typeQNs := termdec.CollectEnv(maps.Values(procDecs))
		typeDefs, getErr3 := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (map[uniqsym.ADT]typedef.DefRec, error) {
			return s.typeDefRepo.SelectEnv(txCtx, conn, typeQNs)
		})
		if getErr3 != nil {
			s.log.Error("step taking failed", execAttr, slog.Any("types", typeQNs))
			return getErr3
		}
		envIDs := typedef.CollectEnv(maps.Values(typeDefs))
		ctxIDs := CollectCtx(maps.Values(evalSnap.LinearVars))
		typeExps, getErr4 := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (map[contkey.ADT]typedef.ExpRec, error) {
			return s.typeExpRepo.SelectEnv(txCtx, conn, append(envIDs, ctxIDs...))
		})
		if getErr4 != nil {
			s.log.Error("step taking failed", execAttr, slog.Any("env", envIDs), slog.Any("ctx", ctxIDs))
			return getErr4
		}
		procEnv := Env{ProcDecs: procDecs, TypeDefs: typeDefs, TypeExps: typeExps}
		procCtx := convertToCtx(maps.Values(evalSnap.LinearVars), typeExps)
		// type checking
		err = s.checkType(procEnv, procCtx, evalSnap, expSpec)
		if err != nil {
			s.log.Error("step taking failed", execAttr)
			return err
		}
		// step taking
		result, err := s.takeWith(ctx, procEnv, evalSnap, expSpec)
		if err != nil {
			s.log.Error("step taking failed", execAttr)
			return err
		}
		if len(result.Communication.Acts) > 0 || (result.Communication.OffsetNr != nil && !result.Communication.OffsetNr.IsEmpty()) {
			providerExecRef := execroot.RootRef{ExecID: result.Communication.ProviderExecID, ExecRN: result.Communication.ProviderExecRN}
			if providerExecRef.ExecID == evalSnap.RootRef.ExecID {
				if providerExecRef.ExecRN != evalSnap.RootRef.ExecRN {
					return fmt.Errorf("provider execution root snapshot changed: expected revision %v, got %v", providerExecRef.ExecRN, evalSnap.RootRef.ExecRN)
				}
			} else {
				result.Mod.ExecRefs = append(result.Mod.ExecRefs, providerExecRef)
			}
		}
		err = s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			if len(result.Communication.Acts) > 0 ||
				(result.Communication.OffsetNr != nil && !result.Communication.OffsetNr.IsEmpty()) {
				if err := s.commRepo.ModifyRec(txCtx, conn, result.Communication); err != nil {
					s.log.Error("Communication modification failed", execAttr)
					return err
				}
			}
			if err := s.execEvalRepo.ModifyRec(txCtx, conn, result.Mod); err != nil {
				s.log.Error("step taking failed", execAttr)
				return err
			}
			return nil
		})
		if err != nil {
			s.log.Error("step taking failed", execAttr)
			return err
		}
		// next values
		if len(result.Effects.Steps) > 0 {
			nextSpec := result.Effects.Steps[len(result.Effects.Steps)-1]
			execRef = nextSpec.RootRef
			expSpec = nextSpec.ProcExp
		}
	}
	s.log.Debug("step taking succeed", execAttr)
	return nil
}

func (s *service) takeWith(
	ctx context.Context,
	procEnv Env,
	evalSnap EvalSnap,
	exp commact.ExpSpec,
) (takeResult, error) {
	var evalMod EvalMod
	var evalEff EvalEff
	var commMod execcomm.CommMod
	var err error
	execAttr := slog.Any("rootRef", evalSnap.RootRef)
	evalMod.ExecRefs = append(evalMod.ExecRefs, evalSnap.RootRef)
	switch termExp := exp.(type) {
	case commact.CloseSpec:
		commChnl, ok := evalSnap.LinearVars[termExp.ContChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, termdef.ErrMissingInCfg(termExp.ContChnlPH)
		}
		commAttr := slog.Any("commRef", commChnl.CommID)
		// получаем снепшот коммуникации
		var commSnap execcomm.CommSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			commSnap, err = s.commRepo.GetSnapByQry(txCtx, conn, execcomm.CommQry{
				CommID: commChnl.CommID,
				ChnlID: option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", execAttr, commAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, getErr
		}
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		subscription := commSnap.NextAct()
		if subscription != nil {
			commMod.CommID = commSnap.CommID
			commMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку закрывателя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: commact.CloseRec(termExp),
			})
			s.log.Debug("taking half done", execAttr, commAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		}
		observation, ok := subscription.(commact.SubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		nextExpCK := contkey.One.Invert()
		switch expRec := observation.ContExp.(type) {
		case commact.WaitRec:
			// лишаем продолжения закрывателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: commChnl.RootRef,
				CommID:  commChnl.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// лишаем продолжения наблюдателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: observation.RootRef,
				CommID:  observation.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение наблюдателя
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: observation.RootRef,
					ProcExp: expRec.ContExp,
				})
			}
			s.log.Debug("step taking succeed", execAttr, commAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		default:
			panic(commact.ErrRecTypeUnexpected(observation.ContExp))
		}
	case commact.WaitSpec:
		commChnl, ok := evalSnap.LinearVars[termExp.ContChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, termdef.ErrMissingInCfg(termExp.ContChnlPH)
		}
		// получаем снепшот соединения
		var commSnap execcomm.CommSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			commSnap, err = s.commRepo.GetSnapByQry(txCtx, conn, execcomm.CommQry{
				CommID: commChnl.CommID,
				ChnlID: option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, getErr
		}
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		subscription := commSnap.NextAct()
		if subscription != nil {
			commMod.CommID = commSnap.CommID
			commMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку наблюдателя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: commact.WaitRec(termExp),
			})
			s.log.Debug("taking half done", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		}
		closage, ok := subscription.(commact.SubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		nextExpCK := contkey.One.Invert()
		switch expRec := closage.ContExp.(type) {
		case commact.CloseRec:
			// лишаем продолжения наблюдателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: commChnl.RootRef,
				CommID:  commChnl.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// лишаем продолжения закрывателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: closage.RootRef,
				CommID:  closage.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if termExp.ContExp != nil {
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: evalSnap.RootRef,
					ProcExp: termExp.ContExp,
				})
			}
			s.log.Debug("step taking succeed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		case commact.FwdRec:
			// перенаправляем продолжение наблюдателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: evalSnap.RootRef,
				CommID:  commSnap.CommID,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpCK:   commChnl.ExpCK,
			})
			evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
				RootRef: evalSnap.RootRef,
				ProcExp: termExp,
			})
			s.log.Debug("step taking succeed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		default:
			panic(commact.ErrRecTypeUnexpected(closage.ContExp))
		}
	case commact.SendSpec:
		commChnl, ok := evalSnap.LinearVars[termExp.CommChnlPH]
		if !ok {
			err := termdef.ErrMissingInCfg(termExp.CommChnlPH)
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, err
		}
		typeExp, ok := procEnv.TypeExps[commChnl.ExpCK]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, typedef.ErrMissingInEnv(commChnl.ExpCK)
		}
		nextExpCK := typeExp.(typedef.ProdRec).Next()
		valChnl, ok := evalSnap.LinearVars[termExp.ValChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, termdef.ErrMissingInCfg(termExp.ValChnlPH)
		}
		// получаем снепшот соединения
		var commSnap execcomm.CommSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			commSnap, err = s.commRepo.GetSnapByQry(txCtx, conn, execcomm.CommQry{
				CommID: commChnl.CommID,
				ChnlID: option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, getErr
		}
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		subscription := commSnap.NextAct()
		if subscription != nil {
			commMod.CommID = commSnap.CommID
			commMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку отправителя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: commact.SendRec{
					CommChnlPH: commChnl.ChnlPH,
					ValChnlID:  valChnl.ChnlID,
					ValExpCK:   valChnl.ExpCK,
				},
			})
			s.log.Debug("taking half done", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		}
		receival, ok := subscription.(commact.SubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение отправителя
		evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
			RootRef: commChnl.RootRef,
			CommID:  commChnl.CommID,
			ChnlID:  newChnlID,
			ChnlPH:  commChnl.ChnlPH,
			ChnlBS:  commChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		// лишаем значения отправителя
		evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
			RootRef: commChnl.RootRef,
			CommID:  commChnl.CommID,
			ChnlID:  identity.New(),
			ChnlPH:  valChnl.ChnlPH,
			ChnlBS:  valChnl.ChnlBS,
			ExpCK:   valChnl.ExpCK.Invert(),
		})
		switch expRec := receival.ContExp.(type) {
		case commact.RecvRec:
			// вяжем продолжение принимателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: receival.RootRef,
				CommID:  receival.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.CommChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// вяжем значение принимателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: receival.RootRef,
				CommID:  valChnl.CommID,
				ChnlID:  valChnl.ChnlID,
				ChnlPH:  expRec.NewChnlPH,
				ChnlBS:  valChnl.ChnlBS,
				ExpCK:   valChnl.ExpCK,
			})
			if expRec.ContExp != nil {
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: receival.RootRef,
					ProcExp: expRec.ContExp,
				})
			}
			s.log.Debug("step taking succeed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		default:
			panic(commact.ErrRecTypeUnexpected(receival.ContExp))
		}
	case commact.RecvSpec:
		commChnl, ok := evalSnap.LinearVars[termExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, termdef.ErrMissingInCfg(termExp.CommChnlPH)
		}
		typeExp, ok := procEnv.TypeExps[commChnl.ExpCK]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, typedef.ErrMissingInEnv(commChnl.ExpCK)
		}
		nextExpCK := typeExp.(typedef.ProdRec).Next()
		// получаем снепшот соединения
		var commSnap execcomm.CommSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			commSnap, err = s.commRepo.GetSnapByQry(txCtx, conn, execcomm.CommQry{
				CommID: commChnl.CommID,
				ChnlID: option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, getErr
		}
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		publication := commSnap.NextAct()
		if publication != nil {
			commMod.CommID = commSnap.CommID
			commMod.OffsetNr = option.Some(publication.GetActNr())
		}
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение принимателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: commChnl.RootRef,
				CommID:  commChnl.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// регистрируем подписку принимателя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: commact.RecvRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					NewChnlPH:  termExp.NewChnlPH,
					ContExp:    termExp.ContExp,
				},
			})
			s.log.Debug("taking half done", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		}
		sending, ok := publication.(commact.PubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(publication))
		}
		switch valExp := sending.ValExp.(type) {
		case commact.SendRec:
			// вяжем продолжение принимателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: commChnl.RootRef,
				CommID:  commChnl.CommID,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// вяжем значение принимателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: commChnl.RootRef,
				CommID:  valExp.CommID,
				ChnlID:  valExp.ValChnlID,
				ChnlPH:  termExp.NewChnlPH,
				ChnlBS:  execvar.AssetSide,
				ExpCK:   valExp.ValExpCK,
			})
			if termExp.ContExp != nil {
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: evalSnap.RootRef,
					ProcExp: termExp.ContExp,
				})
			}
			s.log.Debug("step taking succeed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		default:
			panic(commact.ErrRecTypeUnexpected(sending.ValExp))
		}
	case commact.LabSpec:
		commChnl, ok := evalSnap.LinearVars[termExp.CommChnlPH]
		if !ok {
			err := termdef.ErrMissingInCfg(termExp.CommChnlPH)
			s.log.Error("step taking failed")
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, err
		}
		typeExp, ok := procEnv.TypeExps[commChnl.ExpCK]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, typedef.ErrMissingInEnv(commChnl.ExpCK)
		}
		nextExpCK := typeExp.(typedef.SumRec).Next(termExp.ValLabQN)
		// получаем снепшот соединения
		var commSnap execcomm.CommSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			commSnap, err = s.commRepo.GetSnapByQry(txCtx, conn, execcomm.CommQry{
				CommID: commChnl.CommID,
				ChnlID: option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, getErr
		}
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		subscription := commSnap.NextAct()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение решателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: commChnl.RootRef,
				CommID:  commChnl.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// регистрируем сообщение решателя
			commMod.Acts = append(commMod.Acts, commact.PubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ValExp: commact.LabRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					ValLabQN:   termExp.ValLabQN,
				},
			})
			s.log.Debug("taking half done", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		}
		following, ok := subscription.(commact.SubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(subscription))
		}
		switch contExp := following.ContExp.(type) {
		case commact.CaseRec:
			// вяжем продолжение решателя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: commChnl.RootRef,
				CommID:  commChnl.CommID,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// вяжем продолжение последователя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: following.RootRef,
				CommID:  following.CommID,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  contExp.CommChnlPH,
				// TODO значение ChnlBS
				ExpCK: nextExpCK,
			})
			if contExp.ContExps[termExp.ValLabQN] != nil {
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: following.RootRef,
					ProcExp: contExp.ContExps[termExp.ValLabQN],
				})
			}
			s.log.Debug("step taking succeed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		default:
			panic(commact.ErrRecTypeUnexpected(following.ContExp))
		}
	case commact.CaseSpec:
		commChnl, ok := evalSnap.LinearVars[termExp.CommChnlPH]
		if !ok {
			err := termdef.ErrMissingInCfg(termExp.CommChnlPH)
			s.log.Error("step taking failed")
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, err
		}
		// получаем снепшот соединения
		var caseCommSnap execcomm.CommSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			caseCommSnap, err = s.commRepo.GetSnapByQry(txCtx, conn, execcomm.CommQry{
				CommID: commChnl.CommID,
				ChnlID: option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, getErr
		}
		commMod.ProviderExecID = caseCommSnap.ProviderExecID
		commMod.ProviderExecRN = caseCommSnap.ProviderExecRN
		publication := caseCommSnap.NextAct()
		if publication != nil {
			commMod.CommID = caseCommSnap.CommID
			commMod.OffsetNr = option.Some(publication.GetActNr())
		}
		if publication == nil {
			newChnlID := identity.New()
			// регистрируем подписку последователя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commChnl.CommID,
				RootRef: commChnl.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExp: commact.CaseRec{
					CommChnlPH: commChnl.ChnlPH,
					ContChnlID: newChnlID,
					ContExps:   termExp.ContExps,
				},
			})
			s.log.Debug("taking half done", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		}
		decision, ok := publication.(commact.PubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(publication))
		}
		switch valExp := decision.ValExp.(type) {
		case commact.LabRec:
			// вяжем продолжение последователя
			evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
				RootRef: commChnl.RootRef,
				CommID:  commChnl.CommID,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  commChnl.ChnlPH,
				// TODO значение ChnlBS
				ExpCK: valExp.ValExpCK,
			})
			if termExp.ContExps[valExp.ValLabQN] != nil {
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: evalSnap.RootRef,
					ProcExp: termExp.ContExps[valExp.ValLabQN],
				})
			}
			s.log.Debug("step taking succeed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
		default:
			panic(commact.ErrRecTypeUnexpected(decision.ValExp))
		}
	case commact.FwdSpec:
		commChnl, ok := evalSnap.LinearVars[termExp.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed")
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, termdef.ErrMissingInCfg(termExp.CommChnlPH)
		}
		contChnl, ok := evalSnap.LinearVars[termExp.ContChnlPH]
		if !ok {
			s.log.Error("step taking failed")
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, termdef.ErrMissingInCfg(termExp.ContChnlPH)
		}
		typeExp, ok := procEnv.TypeExps[commChnl.ExpCK]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, typedef.ErrMissingInEnv(commChnl.ExpCK)
		}
		// получаем снепшот соединения
		var fwdCommSnap execcomm.CommSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			fwdCommSnap, err = s.commRepo.GetSnapByQry(txCtx, conn, execcomm.CommQry{
				CommID: commChnl.CommID,
				ChnlID: option.Some(commChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", execAttr)
			return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, getErr
		}
		commMod.ProviderExecID = fwdCommSnap.ProviderExecID
		commMod.ProviderExecRN = fwdCommSnap.ProviderExecRN
		act := fwdCommSnap.NextAct()
		if act != nil {
			commMod.CommID = fwdCommSnap.CommID
			commMod.OffsetNr = option.Some(act.GetActNr())
		}
		switch typeExp.Pol() {
		case polarity.Pos:
			switch forwardable := act.(type) {
			case commact.SubRec:
				// перенаправляем подписчика
				evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
					RootRef: forwardable.RootRef,
					CommID:  forwardable.CommID,
					ChnlID:  commChnl.ChnlID,
					ChnlPH:  forwardable.ContExp.Via(),
					// TODO значение ChnlBS
					ExpCK: commChnl.ExpCK,
				})
				if forwardable.ContExp != nil {
					evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
						RootRef: forwardable.RootRef,
						ProcExp: forwardable.ContExp,
					})
				}
				s.log.Debug("step taking succeed", execAttr)
				return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
			case commact.PubRec:
				// перенаправляем публикатора
				evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
					RootRef: forwardable.RootRef,
					CommID:  forwardable.CommID,
					ChnlID:  contChnl.ChnlID,
					ChnlPH:  forwardable.ValExp.Via(),
					// TODO значение ChnlBS
					ExpCK: contChnl.ExpCK,
				})
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: forwardable.RootRef,
					ProcExp: forwardable.ValExp,
				})
				s.log.Debug("step taking succeed", execAttr)
				return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
			case nil:
				// лишаем значений схлопывающегося
				evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
					RootRef: commChnl.RootRef,
					CommID:  commChnl.CommID,
					ChnlID:  commChnl.ChnlID,
					ChnlPH:  termExp.CommChnlPH,
					ChnlBS:  commChnl.ChnlBS,
					ExpCK:   commChnl.ExpCK.Invert(),
				})
				evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
					RootRef: commChnl.RootRef,
					// TODO значение CommID
					ChnlID: contChnl.ChnlID,
					ChnlPH: termExp.ContChnlPH,
					ChnlBS: contChnl.ChnlBS,
					ExpCK:  contChnl.ExpCK.Invert(),
				})
				// регистрируем сообщение схлопывающегося
				commMod.Acts = append(commMod.Acts, commact.PubRec{
					CommID:  fwdCommSnap.CommID,
					RootRef: evalSnap.RootRef,
					ChnlID:  commChnl.ChnlID,
					ValExp: commact.FwdRec{
						ContChnlID: contChnl.ChnlID,
					},
				})
				s.log.Debug("taking half done", execAttr)
				return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
			default:
				panic(commact.ErrActRecTypeUnexpected(act))
			}
		case polarity.Neg:
			switch forwardable := act.(type) {
			case commact.SubRec:
				// перенаправляем подписчика
				evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
					RootRef: forwardable.RootRef,
					CommID:  forwardable.CommID,
					ChnlID:  contChnl.ChnlID,
					ChnlPH:  forwardable.ContExp.Via(),
					// TODO значение ChnlBS
					ExpCK: contChnl.ExpCK,
				})
				if forwardable.ContExp != nil {
					evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
						RootRef: forwardable.RootRef,
						ProcExp: forwardable.ContExp,
					})
				}
				s.log.Debug("step taking succeed", execAttr)
				return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
			case commact.PubRec:
				// перенаправляем публикатора
				evalMod.LinearVars = append(evalMod.LinearVars, execvar.LinearRec{
					RootRef: forwardable.RootRef,
					CommID:  forwardable.CommID,
					ChnlID:  commChnl.ChnlID,
					ChnlPH:  forwardable.ValExp.Via(),
					// TODO значение ChnlBS
					ExpCK: commChnl.ExpCK,
				})
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: forwardable.RootRef,
					ProcExp: forwardable.ValExp, // TODO: несовпадение типов
				})
				s.log.Debug("step taking succeed", execAttr)
				return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
			case nil:
				// регистрируем подписку схлопывающегося
				commMod.Acts = append(commMod.Acts, commact.SubRec{
					CommID:  fwdCommSnap.CommID,
					RootRef: commChnl.RootRef,
					ChnlID:  commChnl.ChnlID,
					ContExp: commact.FwdRec{
						ContChnlID: contChnl.ChnlID,
					},
				})
				s.log.Debug("taking half done", execAttr)
				return takeResult{Mod: evalMod, Effects: evalEff, Communication: commMod}, nil
			default:
				panic(commact.ErrActRecTypeUnexpected(act))
			}
		default:
			panic(typedef.ErrPolarityUnexpected(typeExp))
		}
	default:
		panic(commact.ErrExpTypeUnexpected(exp))
	}
}

func CollectCtx(chnls iter.Seq[execvar.LinearRec]) []contkey.ADT {
	return nil
}

func convertToCtx(chnlBinds iter.Seq[execvar.LinearRec], typeExps map[contkey.ADT]typedef.ExpRec) typedef.Context {
	assets := make(map[symbol.ADT]typedef.ExpRec, 1)
	liabs := make(map[symbol.ADT]typedef.ExpRec, 1)
	for bind := range chnlBinds {
		if bind.ChnlBS == execvar.LiabSide {
			liabs[bind.ChnlPH] = typeExps[bind.ExpCK]
		} else {
			assets[bind.ChnlPH] = typeExps[bind.ExpCK]
		}
	}
	return typedef.Context{Assets: assets, Liabs: liabs}
}

func (s *service) checkType(
	procEnv Env,
	procCtx typedef.Context,
	evalSnap EvalSnap,
	expSpec commact.ExpSpec,
) error {
	chnlBR, ok := evalSnap.LinearVars[expSpec.Via()]
	if !ok {
		panic("no comm chnl in proc snap")
	}
	if chnlBR.ChnlBS == execvar.LiabSide {
		return s.checkProvider(procEnv, procCtx, evalSnap, expSpec)
	}
	return s.checkClient(procEnv, procCtx, evalSnap, expSpec)
}

func (s *service) checkProvider(
	procEnv Env,
	procCtx typedef.Context,
	procCfg EvalSnap,
	es commact.ExpSpec,
) error {
	switch expSpec := es.(type) {
	case commact.CloseSpec:
		// check ctx
		if len(procCtx.Assets) > 0 {
			err := fmt.Errorf("context mismatch: want 0 items, got %v items", len(procCtx.Assets))
			s.log.Error("checking failed")
			return err
		}
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.ContChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.ContChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(gotVia, typedef.OneRec{})
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// no cont to check
		delete(procCtx.Liabs, expSpec.ContChnlPH)
		return nil
	case commact.WaitSpec:
		err := commact.ErrExpTypeMismatch(es, commact.CloseSpec{})
		s.log.Error("checking failed")
		return err
	case commact.SendSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.CommChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.CommChnlPH)
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
		gotVal, ok := procCtx.Assets[expSpec.ValChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.ValChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(gotVal, wantVia.Val)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// no cont to check
		procCtx.Liabs[expSpec.CommChnlPH] = wantVia.Cont
		delete(procCtx.Assets, expSpec.ValChnlPH)
		return nil
	case commact.RecvSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.CommChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.CommChnlPH)
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
		gotVal, ok := procCtx.Assets[expSpec.NewChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.NewChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(gotVal, wantVia.Val)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// check cont
		procCtx.Liabs[expSpec.CommChnlPH] = wantVia.Cont
		procCtx.Assets[expSpec.NewChnlPH] = wantVia.Val
		return s.checkType(procEnv, procCtx, procCfg, expSpec.ContExp)
	case commact.LabSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.CommChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.CommChnlPH)
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
		choice, ok := wantVia.Choices[expSpec.ValLabQN]
		if !ok {
			err := fmt.Errorf("label mismatch: want %v, got %v", maps.Keys(wantVia.Choices), expSpec.ValLabQN)
			s.log.Error("checking failed")
			return err
		}
		// no cont to check
		procCtx.Liabs[expSpec.CommChnlPH] = choice
		return nil
	case commact.CaseSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.CommChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.CommChnlPH)
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
		if len(expSpec.ContExps) != len(wantVia.Choices) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantVia.Choices), len(expSpec.ContExps))
			s.log.Error("checking failed")
			return err
		}
		for label, choice := range wantVia.Choices {
			cont, ok := expSpec.ContExps[label]
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
	case commact.FwdSpec:
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
			err := termdef.ErrMissingInCtx(expSpec.ContChnlPH)
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
		delete(procCtx.Liabs, expSpec.CommChnlPH)
		delete(procCtx.Assets, expSpec.ContChnlPH)
		return nil
	default:
		panic(commact.ErrExpTypeUnexpected(es))
	}
}

func (s *service) checkClient(
	procEnv Env,
	procCtx typedef.Context,
	procCfg EvalSnap,
	es commact.ExpSpec,
) error {
	switch expSpec := es.(type) {
	case commact.CloseSpec:
		err := commact.ErrExpTypeMismatch(es, commact.WaitSpec{})
		s.log.Error("checking failed")
		return err
	case commact.WaitSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.ContChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.ContChnlPH)
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
		delete(procCtx.Assets, expSpec.ContChnlPH)
		return s.checkType(procEnv, procCtx, procCfg, expSpec.ContExp)
	case commact.SendSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.CommChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.CommChnlPH)
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
		gotVal, ok := procCtx.Assets[expSpec.ValChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.ValChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(gotVal, wantVia.Val)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		procCtx.Assets[expSpec.CommChnlPH] = wantVia.Cont
		delete(procCtx.Assets, expSpec.ValChnlPH)
		return nil
	case commact.RecvSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.CommChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.CommChnlPH)
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
		gotVal, ok := procCtx.Assets[expSpec.NewChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.NewChnlPH)
			s.log.Error("checking failed")
			return err
		}
		err := typedef.CheckRec(gotVal, wantVia.Val)
		if err != nil {
			s.log.Error("checking failed")
			return err
		}
		// check cont
		procCtx.Assets[expSpec.CommChnlPH] = wantVia.Cont
		procCtx.Assets[expSpec.NewChnlPH] = wantVia.Val
		return s.checkType(procEnv, procCtx, procCfg, expSpec.ContExp)
	case commact.LabSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.CommChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.CommChnlPH)
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
		choice, ok := wantVia.Choices[expSpec.ValLabQN]
		if !ok {
			err := fmt.Errorf("label mismatch: want %v, got %v", maps.Keys(wantVia.Choices), expSpec.ValLabQN)
			s.log.Error("checking failed")
			return err
		}
		procCtx.Assets[expSpec.CommChnlPH] = choice
		return nil
	case commact.CaseSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.CommChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.CommChnlPH)
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
		if len(expSpec.ContExps) != len(wantVia.Choices) {
			err := fmt.Errorf("state mismatch: want %v choices, got %v conts", len(wantVia.Choices), len(expSpec.ContExps))
			s.log.Error("checking failed")
			return err
		}
		for label, choice := range wantVia.Choices {
			cont, ok := expSpec.ContExps[label]
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
	// 		wantVal, ok := procEnv.TypeExps[valType.ExpCK]
	// 		if !ok {
	// 			err := typedef.ErrMissingInEnv(valType.ExpCK)
	// 			s.log.Error("checking failed")
	// 			return err
	// 		}
	// 		gotVal, ok := procCtx.Assets[expSpec.Ys[i]]
	// 		if !ok {
	// 			err := procdef.ErrMissingInCtx(ep.ChnlPH)
	// 			s.log.Error("checking failed")
	// 			return err
	// 		}
	// 		err := typedef.CheckRec(gotVal, wantVal)
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
	// 	wantVia, ok := procEnv.TypeExps[viaRole.ExpCK]
	// 	if !ok {
	// 		err := typedef.ErrMissingInEnv(viaRole.ExpCK)
	// 		s.log.Error("checking failed")
	// 		return err
	// 	}
	// 	// check cont
	// 	procCtx.Assets[expSpec.X] = wantVia
	// 	return s.checkType(procEnv, procCtx, procCfg, expSpec.ContES)
	default:
		panic(commact.ErrExpTypeUnexpected(es))
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
