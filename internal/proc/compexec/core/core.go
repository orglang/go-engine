package core

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"maps"
	"reflect"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/compvar"
	"orglang/go-engine/internal/adt/contkey"
	compexp "orglang/go-engine/internal/proc/compexec"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/polarity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	compconn "orglang/go-engine/internal/proc/compconn/core"
	compstep "orglang/go-engine/internal/proc/compstep/core"
	connact "orglang/go-engine/internal/proc/connact/core"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	termdef "orglang/go-engine/internal/proc/termdef/core"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type API interface {
	Take(context.Context, compstep.StepSpec) error
	RetrieveSnap(comproot.SemRef) (ExecSnap, error)
}

type ExecRec struct {
	CompRef  comproot.SemRef
	LiabMode compvar.Mode
}

type ExecMod struct {
	CompRefs   []comproot.SemRef
	LinearVars []compvar.LinearRec
}

type ExecEff struct {
	Steps []compstep.StepSpec
}

type takeResult struct {
	Mod        ExecMod
	Effects    ExecEff
	Connection compconn.ConnMod
}

// aka Configuration
type ExecSnap struct {
	CompRef    comproot.SemRef
	LinearVars map[symbol.ADT]compvar.LinearRec
}

type Env struct {
	TypeDefs map[uniqsym.ADT]typedef.DefRec
	TypeExps map[contkey.ADT]typedef.ExpRec
	ProcDecs map[identity.ADT]termdec.DecRec
}

func ChnlPH(rec compvar.LinearRec) symbol.ADT { return rec.ChnlPH }

type service struct {
	compExecRepo Repo
	connRepo     compconn.Repo
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
	compExecRepo Repo,
	connRepo compconn.Repo,
	termDecRepo termdec.Repo,
	typeDefRepo typedef.Repo,
	typeExpRepo typedef.TypeExpRepo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{
		compExecRepo, connRepo,
		termDecRepo, typeDefRepo, typeExpRepo,
		transactor, log.With(name),
	}
}

func (s *service) RetrieveSnap(ref comproot.SemRef) (_ ExecSnap, err error) {
	return ExecSnap{}, nil
}

func ErrMissingChnl(want symbol.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
}

func (s *service) Take(ctx context.Context, spec compstep.StepSpec) (err error) {
	compAttr := slog.Any("proc", spec.CompRef)
	s.log.Debug("step taking started", compAttr, slog.Any("exp", spec.ProcExp))
	// initial values
	compRef := spec.CompRef
	expSpec := spec.ProcExp
	for expSpec != nil {
		execSnap, getErr1 := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (ExecSnap, error) {
			return s.compExecRepo.GetSnapByRef(txCtx, conn, compRef)
		})
		if getErr1 != nil {
			s.log.Error("step taking failed", compAttr)
			return getErr1
		}
		if len(execSnap.LinearVars) == 0 {
			panic("zero channel binds")
		}
		decIDs := compexp.CollectEnv(expSpec)
		procDecs, getErr2 := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (map[identity.ADT]termdec.DecRec, error) {
			return s.termDecRepo.SelectEnv(txCtx, conn, decIDs)
		})
		if getErr2 != nil {
			s.log.Error("step taking failed", compAttr, slog.Any("decs", decIDs))
			return getErr2
		}
		typeQNs := termdec.CollectEnv(maps.Values(procDecs))
		typeDefs, getErr3 := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (map[uniqsym.ADT]typedef.DefRec, error) {
			return s.typeDefRepo.SelectEnv(txCtx, conn, typeQNs)
		})
		if getErr3 != nil {
			s.log.Error("step taking failed", compAttr, slog.Any("types", typeQNs))
			return getErr3
		}
		envIDs := typedef.CollectEnv(maps.Values(typeDefs))
		ctxIDs := CollectCtx(maps.Values(execSnap.LinearVars))
		typeExps, getErr4 := tx.Implicit(ctx, s.transactor, func(txCtx context.Context, conn tx.Conn) (map[contkey.ADT]typedef.ExpRec, error) {
			return s.typeExpRepo.SelectEnv(txCtx, conn, append(envIDs, ctxIDs...))
		})
		if getErr4 != nil {
			s.log.Error("step taking failed", compAttr, slog.Any("env", envIDs), slog.Any("ctx", ctxIDs))
			return getErr4
		}
		procEnv := Env{ProcDecs: procDecs, TypeDefs: typeDefs, TypeExps: typeExps}
		procCtx := convertToCtx(maps.Values(execSnap.LinearVars), typeExps)
		// type checking
		err = s.checkType(procEnv, procCtx, execSnap, expSpec)
		if err != nil {
			s.log.Error("step taking failed", compAttr)
			return err
		}
		// step taking
		result, err := s.takeWith(ctx, procEnv, execSnap, expSpec)
		if err != nil {
			s.log.Error("step taking failed", compAttr)
			return err
		}
		if len(result.Connection.Acts) > 0 || (result.Connection.OffsetNr != nil && !result.Connection.OffsetNr.IsEmpty()) {
			providerRef := comproot.SemRef{CompID: result.Connection.ProviderCompID, CompRN: result.Connection.ProviderCompRN}
			if providerRef.CompID == execSnap.CompRef.CompID {
				if providerRef.CompRN != execSnap.CompRef.CompRN {
					return fmt.Errorf("provider computation snapshot changed: expected revision %v, got %v", providerRef.CompRN, execSnap.CompRef.CompRN)
				}
			} else {
				result.Mod.CompRefs = append(result.Mod.CompRefs, providerRef)
			}
		}
		err = s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			if len(result.Connection.Acts) > 0 ||
				(result.Connection.OffsetNr != nil && !result.Connection.OffsetNr.IsEmpty()) {
				if err := s.connRepo.ModifyRec(txCtx, conn, result.Connection); err != nil {
					s.log.Error("Connection modification failed", compAttr)
					return err
				}
			}
			if err := s.compExecRepo.ModifyRec(txCtx, conn, result.Mod); err != nil {
				s.log.Error("step taking failed", compAttr)
				return err
			}
			return nil
		})
		if err != nil {
			s.log.Error("step taking failed", compAttr)
			return err
		}
		// next values
		if len(result.Effects.Steps) > 0 {
			nextSpec := result.Effects.Steps[len(result.Effects.Steps)-1]
			compRef = nextSpec.CompRef
			expSpec = nextSpec.ProcExp
		}
	}
	s.log.Debug("step taking succeed", compAttr)
	return nil
}

func (s *service) takeWith(
	ctx context.Context,
	procEnv Env,
	execSnap ExecSnap,
	exp compexp.ExpSpec,
) (takeResult, error) {
	var execMod ExecMod
	var execEff ExecEff
	var connMod compconn.ConnMod
	var err error
	compAttr := slog.Any("compRef", execSnap.CompRef)
	execMod.CompRefs = append(execMod.CompRefs, execSnap.CompRef)
	switch termExp := exp.(type) {
	case compexp.CloseSpec:
		connChnl, ok := execSnap.LinearVars[termExp.ContChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, termdef.ErrMissingInCfg(termExp.ContChnlPH)
		}
		connAttr := slog.Any("connRef", connChnl.ConnID)
		// получаем снепшот коммуникации
		var connSnap compconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr, connAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, getErr
		}
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.ConnID = connSnap.ConnID
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку закрывателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: compexp.CloseRec(termExp),
			})
			s.log.Debug("taking half done", compAttr, connAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		}
		observation, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		nextExpCK := contkey.One.Invert()
		switch expRec := observation.ContExp.(type) {
		case compexp.WaitRec:
			// лишаем продолжения закрывателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: connChnl.CompRef,
				ConnID:  connChnl.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  connChnl.ChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// лишаем продолжения наблюдателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: observation.CompRef,
				ConnID:  observation.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение наблюдателя
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: observation.CompRef,
					ProcExp: expRec.ContExp,
				})
			}
			s.log.Debug("step taking succeed", compAttr, connAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		default:
			panic(compexp.ErrRecTypeUnexpected(observation.ContExp))
		}
	case compexp.WaitSpec:
		connChnl, ok := execSnap.LinearVars[termExp.ContChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, termdef.ErrMissingInCfg(termExp.ContChnlPH)
		}
		// получаем снепшот соединения
		var connSnap compconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, getErr
		}
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.ConnID = connSnap.ConnID
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку наблюдателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: compexp.WaitRec(termExp),
			})
			s.log.Debug("taking half done", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		}
		closage, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		nextExpCK := contkey.One.Invert()
		switch expRec := closage.ContExp.(type) {
		case compexp.CloseRec:
			// лишаем продолжения наблюдателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: connChnl.CompRef,
				ConnID:  connChnl.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  connChnl.ChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// лишаем продолжения закрывателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: closage.CompRef,
				ConnID:  closage.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if termExp.ContExp != nil {
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: execSnap.CompRef,
					ProcExp: termExp.ContExp,
				})
			}
			s.log.Debug("step taking succeed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		case compexp.FwdRec:
			// перенаправляем продолжение наблюдателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: execSnap.CompRef,
				ConnID:  connSnap.ConnID,
				ChnlID:  expRec.ContChnlID,
				ChnlPH:  connChnl.ChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   connChnl.ExpCK,
			})
			execEff.Steps = append(execEff.Steps, compstep.StepSpec{
				CompRef: execSnap.CompRef,
				ProcExp: termExp,
			})
			s.log.Debug("step taking succeed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		default:
			panic(compexp.ErrRecTypeUnexpected(closage.ContExp))
		}
	case compexp.SendSpec:
		connChnl, ok := execSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			err := termdef.ErrMissingInCfg(termExp.ConnChnlPH)
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, err
		}
		typeExp, ok := procEnv.TypeExps[connChnl.ExpCK]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, typedef.ErrMissingInEnv(connChnl.ExpCK)
		}
		nextExpCK := typeExp.(typedef.ProdRec).Next()
		valChnl, ok := execSnap.LinearVars[termExp.ValChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, termdef.ErrMissingInCfg(termExp.ValChnlPH)
		}
		// получаем снепшот соединения
		var connSnap compconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, getErr
		}
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.ConnID = connSnap.ConnID
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку отправителя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: compexp.SendRec{
					ConnChnlPH: connChnl.ChnlPH,
					ValChnlID:  valChnl.ChnlID,
					ValExpCK:   valChnl.ExpCK,
				},
			})
			s.log.Debug("taking half done", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		}
		receival, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение отправителя
		execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
			CompRef: connChnl.CompRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		// лишаем значения отправителя
		execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
			CompRef: connChnl.CompRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  identity.New(),
			ChnlPH:  valChnl.ChnlPH,
			ChnlBS:  valChnl.ChnlBS,
			ExpCK:   valChnl.ExpCK.Invert(),
		})
		switch expRec := receival.ContExp.(type) {
		case compexp.RecvRec:
			// вяжем продолжение принимателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: receival.CompRef,
				ConnID:  receival.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ConnChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// вяжем значение принимателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: receival.CompRef,
				ConnID:  valChnl.ConnID,
				ChnlID:  valChnl.ChnlID,
				ChnlPH:  expRec.NewChnlPH,
				ChnlBS:  valChnl.ChnlBS,
				ExpCK:   valChnl.ExpCK,
			})
			if expRec.ContExp != nil {
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: receival.CompRef,
					ProcExp: expRec.ContExp,
				})
			}
			s.log.Debug("step taking succeed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		default:
			panic(compexp.ErrRecTypeUnexpected(receival.ContExp))
		}
	case compexp.RecvSpec:
		connChnl, ok := execSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, termdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		typeExp, ok := procEnv.TypeExps[connChnl.ExpCK]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, typedef.ErrMissingInEnv(connChnl.ExpCK)
		}
		nextExpCK := typeExp.(typedef.ProdRec).Next()
		// получаем снепшот соединения
		var connSnap compconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, getErr
		}
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		publication := connSnap.NextAct()
		if publication != nil {
			connMod.ConnID = connSnap.ConnID
			connMod.OffsetNr = option.Some(publication.GetActNr())
		}
		if publication == nil {
			newChnlID := identity.New()
			// вяжем продолжение принимателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: connChnl.CompRef,
				ConnID:  connChnl.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  connChnl.ChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// регистрируем подписку принимателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: compexp.RecvRec{
					ConnChnlPH: connChnl.ChnlPH,
					ContChnlID: newChnlID,
					NewChnlPH:  termExp.NewChnlPH,
					ContExp:    termExp.ContExp,
				},
			})
			s.log.Debug("taking half done", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		}
		sending, ok := publication.(connact.PubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(publication))
		}
		switch valExp := sending.ValExp.(type) {
		case compexp.SendRec:
			// вяжем продолжение принимателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: connChnl.CompRef,
				ConnID:  connChnl.ConnID,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  connChnl.ChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// вяжем значение принимателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: connChnl.CompRef,
				ConnID:  valExp.ConnID,
				ChnlID:  valExp.ValChnlID,
				ChnlPH:  termExp.NewChnlPH,
				ChnlBS:  compvar.AssetSide,
				ExpCK:   valExp.ValExpCK,
			})
			if termExp.ContExp != nil {
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: execSnap.CompRef,
					ProcExp: termExp.ContExp,
				})
			}
			s.log.Debug("step taking succeed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		default:
			panic(compexp.ErrRecTypeUnexpected(sending.ValExp))
		}
	case compexp.LabSpec:
		connChnl, ok := execSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			err := termdef.ErrMissingInCfg(termExp.ConnChnlPH)
			s.log.Error("step taking failed")
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, err
		}
		typeExp, ok := procEnv.TypeExps[connChnl.ExpCK]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, typedef.ErrMissingInEnv(connChnl.ExpCK)
		}
		nextExpCK := typeExp.(typedef.SumRec).Next(termExp.ValLabQN)
		// получаем снепшот соединения
		var connSnap compconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, getErr
		}
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		subscription := connSnap.NextAct()
		if subscription == nil {
			newChnlID := identity.New()
			// вяжем продолжение решателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: connChnl.CompRef,
				ConnID:  connChnl.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  connChnl.ChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// регистрируем сообщение решателя
			connMod.Acts = append(connMod.Acts, connact.PubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ValExp: compexp.LabRec{
					ConnChnlPH: connChnl.ChnlPH,
					ContChnlID: newChnlID,
					ValLabQN:   termExp.ValLabQN,
				},
			})
			s.log.Debug("taking half done", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		}
		following, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		switch contExp := following.ContExp.(type) {
		case compexp.CaseRec:
			// вяжем продолжение решателя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: connChnl.CompRef,
				ConnID:  connChnl.ConnID,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  connChnl.ChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			// вяжем продолжение последователя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: following.CompRef,
				ConnID:  following.ConnID,
				ChnlID:  contExp.ContChnlID,
				ChnlPH:  contExp.ConnChnlPH,
				// TODO значение ChnlBS
				ExpCK: nextExpCK,
			})
			if contExp.ContExps[termExp.ValLabQN] != nil {
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: following.CompRef,
					ProcExp: contExp.ContExps[termExp.ValLabQN],
				})
			}
			s.log.Debug("step taking succeed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		default:
			panic(compexp.ErrRecTypeUnexpected(following.ContExp))
		}
	case compexp.CaseSpec:
		connChnl, ok := execSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			err := termdef.ErrMissingInCfg(termExp.ConnChnlPH)
			s.log.Error("step taking failed")
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, err
		}
		// получаем снепшот соединения
		var caseConnSnap compconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			caseConnSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, getErr
		}
		connMod.ProviderCompID = caseConnSnap.ProviderCompID
		connMod.ProviderCompRN = caseConnSnap.ProviderCompRN
		publication := caseConnSnap.NextAct()
		if publication != nil {
			connMod.ConnID = caseConnSnap.ConnID
			connMod.OffsetNr = option.Some(publication.GetActNr())
		}
		if publication == nil {
			newChnlID := identity.New()
			// регистрируем подписку последователя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connChnl.ConnID,
				CompRef: connChnl.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: compexp.CaseRec{
					ConnChnlPH: connChnl.ChnlPH,
					ContChnlID: newChnlID,
					ContExps:   termExp.ContExps,
				},
			})
			s.log.Debug("taking half done", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		}
		decision, ok := publication.(connact.PubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(publication))
		}
		switch valExp := decision.ValExp.(type) {
		case compexp.LabRec:
			// вяжем продолжение последователя
			execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
				CompRef: connChnl.CompRef,
				ConnID:  connChnl.ConnID,
				ChnlID:  valExp.ContChnlID,
				ChnlPH:  connChnl.ChnlPH,
				// TODO значение ChnlBS
				ExpCK: valExp.ValExpCK,
			})
			if termExp.ContExps[valExp.ValLabQN] != nil {
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: execSnap.CompRef,
					ProcExp: termExp.ContExps[valExp.ValLabQN],
				})
			}
			s.log.Debug("step taking succeed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
		default:
			panic(compexp.ErrRecTypeUnexpected(decision.ValExp))
		}
	case compexp.FwdSpec:
		connChnl, ok := execSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed")
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, termdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		contChnl, ok := execSnap.LinearVars[termExp.ContChnlPH]
		if !ok {
			s.log.Error("step taking failed")
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, termdef.ErrMissingInCfg(termExp.ContChnlPH)
		}
		typeExp, ok := procEnv.TypeExps[connChnl.ExpCK]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, typedef.ErrMissingInEnv(connChnl.ExpCK)
		}
		// получаем снепшот соединения
		var fwdConnSnap compconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			fwdConnSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, getErr
		}
		connMod.ProviderCompID = fwdConnSnap.ProviderCompID
		connMod.ProviderCompRN = fwdConnSnap.ProviderCompRN
		act := fwdConnSnap.NextAct()
		if act != nil {
			connMod.ConnID = fwdConnSnap.ConnID
			connMod.OffsetNr = option.Some(act.GetActNr())
		}
		switch typeExp.Pol() {
		case polarity.Pos:
			switch forwardable := act.(type) {
			case connact.SubRec:
				// перенаправляем подписчика
				execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
					CompRef: forwardable.CompRef,
					ConnID:  forwardable.ConnID,
					ChnlID:  connChnl.ChnlID,
					ChnlPH:  forwardable.ContExp.Via(),
					// TODO значение ChnlBS
					ExpCK: connChnl.ExpCK,
				})
				if forwardable.ContExp != nil {
					execEff.Steps = append(execEff.Steps, compstep.StepSpec{
						CompRef: forwardable.CompRef,
						ProcExp: forwardable.ContExp,
					})
				}
				s.log.Debug("step taking succeed", compAttr)
				return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
			case connact.PubRec:
				// перенаправляем публикатора
				execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
					CompRef: forwardable.CompRef,
					ConnID:  forwardable.ConnID,
					ChnlID:  contChnl.ChnlID,
					ChnlPH:  forwardable.ValExp.Via(),
					// TODO значение ChnlBS
					ExpCK: contChnl.ExpCK,
				})
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: forwardable.CompRef,
					ProcExp: forwardable.ValExp,
				})
				s.log.Debug("step taking succeed", compAttr)
				return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
			case nil:
				// лишаем значений схлопывающегося
				execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
					CompRef: connChnl.CompRef,
					ConnID:  connChnl.ConnID,
					ChnlID:  connChnl.ChnlID,
					ChnlPH:  termExp.ConnChnlPH,
					ChnlBS:  connChnl.ChnlBS,
					ExpCK:   connChnl.ExpCK.Invert(),
				})
				execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
					CompRef: connChnl.CompRef,
					// TODO значение ConnID
					ChnlID: contChnl.ChnlID,
					ChnlPH: termExp.ContChnlPH,
					ChnlBS: contChnl.ChnlBS,
					ExpCK:  contChnl.ExpCK.Invert(),
				})
				// регистрируем сообщение схлопывающегося
				connMod.Acts = append(connMod.Acts, connact.PubRec{
					ConnID:  fwdConnSnap.ConnID,
					CompRef: execSnap.CompRef,
					ChnlID:  connChnl.ChnlID,
					ValExp: compexp.FwdRec{
						ContChnlID: contChnl.ChnlID,
					},
				})
				s.log.Debug("taking half done", compAttr)
				return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
			default:
				panic(connact.ErrActRecTypeUnexpected(act))
			}
		case polarity.Neg:
			switch forwardable := act.(type) {
			case connact.SubRec:
				// перенаправляем подписчика
				execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
					CompRef: forwardable.CompRef,
					ConnID:  forwardable.ConnID,
					ChnlID:  contChnl.ChnlID,
					ChnlPH:  forwardable.ContExp.Via(),
					// TODO значение ChnlBS
					ExpCK: contChnl.ExpCK,
				})
				if forwardable.ContExp != nil {
					execEff.Steps = append(execEff.Steps, compstep.StepSpec{
						CompRef: forwardable.CompRef,
						ProcExp: forwardable.ContExp,
					})
				}
				s.log.Debug("step taking succeed", compAttr)
				return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
			case connact.PubRec:
				// перенаправляем публикатора
				execMod.LinearVars = append(execMod.LinearVars, compvar.LinearRec{
					CompRef: forwardable.CompRef,
					ConnID:  forwardable.ConnID,
					ChnlID:  connChnl.ChnlID,
					ChnlPH:  forwardable.ValExp.Via(),
					// TODO значение ChnlBS
					ExpCK: connChnl.ExpCK,
				})
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: forwardable.CompRef,
					ProcExp: forwardable.ValExp, // TODO: несовпадение типов
				})
				s.log.Debug("step taking succeed", compAttr)
				return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
			case nil:
				// регистрируем подписку схлопывающегося
				connMod.Acts = append(connMod.Acts, connact.SubRec{
					ConnID:  fwdConnSnap.ConnID,
					CompRef: connChnl.CompRef,
					ChnlID:  connChnl.ChnlID,
					ContExp: compexp.FwdRec{
						ContChnlID: contChnl.ChnlID,
					},
				})
				s.log.Debug("taking half done", compAttr)
				return takeResult{Mod: execMod, Effects: execEff, Connection: connMod}, nil
			default:
				panic(connact.ErrActRecTypeUnexpected(act))
			}
		default:
			panic(typedef.ErrPolarityUnexpected(typeExp))
		}
	default:
		panic(compexp.ErrExpTypeUnexpected(exp))
	}
}

func CollectCtx(chnls iter.Seq[compvar.LinearRec]) []contkey.ADT {
	return nil
}

func convertToCtx(chnlBinds iter.Seq[compvar.LinearRec], typeExps map[contkey.ADT]typedef.ExpRec) typedef.Context {
	assets := make(map[symbol.ADT]typedef.ExpRec, 1)
	liabs := make(map[symbol.ADT]typedef.ExpRec, 1)
	for bind := range chnlBinds {
		if bind.ChnlBS == compvar.LiabSide {
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
	execSnap ExecSnap,
	expSpec compexp.ExpSpec,
) error {
	chnlBR, ok := execSnap.LinearVars[expSpec.Via()]
	if !ok {
		panic("no comm chnl in proc snap")
	}
	if chnlBR.ChnlBS == compvar.LiabSide {
		return s.checkProvider(procEnv, procCtx, execSnap, expSpec)
	}
	return s.checkClient(procEnv, procCtx, execSnap, expSpec)
}

func (s *service) checkProvider(
	procEnv Env,
	procCtx typedef.Context,
	procCfg ExecSnap,
	es compexp.ExpSpec,
) error {
	switch expSpec := es.(type) {
	case compexp.CloseSpec:
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
	case compexp.WaitSpec:
		err := compexp.ErrExpTypeMismatch(es, compexp.CloseSpec{})
		s.log.Error("checking failed")
		return err
	case compexp.SendSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.ConnChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.ConnChnlPH)
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
		procCtx.Liabs[expSpec.ConnChnlPH] = wantVia.Cont
		delete(procCtx.Assets, expSpec.ValChnlPH)
		return nil
	case compexp.RecvSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.ConnChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.ConnChnlPH)
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
		procCtx.Liabs[expSpec.ConnChnlPH] = wantVia.Cont
		procCtx.Assets[expSpec.NewChnlPH] = wantVia.Val
		return s.checkType(procEnv, procCtx, procCfg, expSpec.ContExp)
	case compexp.LabSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.ConnChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.ConnChnlPH)
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
		procCtx.Liabs[expSpec.ConnChnlPH] = choice
		return nil
	case compexp.CaseSpec:
		// check via
		gotVia, ok := procCtx.Liabs[expSpec.ConnChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.ConnChnlPH)
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
			procCtx.Liabs[expSpec.ConnChnlPH] = choice
			err := s.checkType(procEnv, procCtx, procCfg, cont)
			if err != nil {
				s.log.Error("checking failed")
				return err
			}
		}
		return nil
	case compexp.FwdSpec:
		if len(procCtx.Assets) != 1 {
			err := fmt.Errorf("context mismatch: want 1 item, got %v items", len(procCtx.Assets))
			s.log.Error("checking failed")
			return err
		}
		viaSt, ok := procCtx.Liabs[expSpec.ConnChnlPH]
		if !ok {
			err := typedef.ErrMissingInCtx(expSpec.ConnChnlPH)
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
		delete(procCtx.Liabs, expSpec.ConnChnlPH)
		delete(procCtx.Assets, expSpec.ContChnlPH)
		return nil
	default:
		panic(compexp.ErrExpTypeUnexpected(es))
	}
}

func (s *service) checkClient(
	procEnv Env,
	procCtx typedef.Context,
	procCfg ExecSnap,
	es compexp.ExpSpec,
) error {
	switch expSpec := es.(type) {
	case compexp.CloseSpec:
		err := compexp.ErrExpTypeMismatch(es, compexp.WaitSpec{})
		s.log.Error("checking failed")
		return err
	case compexp.WaitSpec:
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
	case compexp.SendSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.ConnChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.ConnChnlPH)
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
		procCtx.Assets[expSpec.ConnChnlPH] = wantVia.Cont
		delete(procCtx.Assets, expSpec.ValChnlPH)
		return nil
	case compexp.RecvSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.ConnChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.ConnChnlPH)
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
		procCtx.Assets[expSpec.ConnChnlPH] = wantVia.Cont
		procCtx.Assets[expSpec.NewChnlPH] = wantVia.Val
		return s.checkType(procEnv, procCtx, procCfg, expSpec.ContExp)
	case compexp.LabSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.ConnChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.ConnChnlPH)
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
		procCtx.Assets[expSpec.ConnChnlPH] = choice
		return nil
	case compexp.CaseSpec:
		// check via
		gotVia, ok := procCtx.Assets[expSpec.ConnChnlPH]
		if !ok {
			err := termdef.ErrMissingInCtx(expSpec.ConnChnlPH)
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
			procCtx.Assets[expSpec.ConnChnlPH] = choice
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
		panic(compexp.ErrExpTypeUnexpected(es))
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
