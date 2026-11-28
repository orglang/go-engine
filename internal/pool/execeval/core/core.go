package core

import (
	"context"
	"log/slog"
	"reflect"
	"slices"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	commact "orglang/go-engine/internal/pool/commact/core"
	execcomm "orglang/go-engine/internal/pool/execcomm/core"
	execstep "orglang/go-engine/internal/pool/execstep/core"
	execvar "orglang/go-engine/internal/pool/execvar/core"
	termdec "orglang/go-engine/internal/pool/termdec/core"
	typedef "orglang/go-engine/internal/pool/typedef/core"
	proctermdef "orglang/go-engine/internal/proc/termdef/core"

	procexeceval "orglang/go-engine/internal/proc/execeval/core"
)

type API interface {
	Run(context.Context, EvalSpec) (execroot.RootRef, error) // aka Create
	Take(context.Context, execstep.StepSpec) error
	Spawn(context.Context, execstep.StepSpec) (execroot.RootRef, error)
}

type Broker interface {
	Subscribe(api API)
	SendSpec(execstep.StepSpec) error
}

type EvalSpec struct {
	// ссылка на декларацию вновь создаваемого пула
	TermQN uniqsym.ADT
	// внутренняя и внешняя ссылки на вновь создаваемый пул
	LiabVar execvar.VarSpec
	// внутренние и внешние ссылки на ранее созданные пулы
	AssetVars []execvar.VarSpec
}

type EvalRec struct {
	RootRef  execroot.RootRef
	TermQN   uniqsym.ADT
	LiabMode execvar.Mode
}

type EvalMod struct {
	RootRef execroot.RootRef
	Vars    []execvar.VarRec
}

func (mod EvalMod) isEmpty() bool { return len(mod.Vars) == 0 }

type EvalEff struct {
	Steps []execstep.StepSpec
}

type EvalSnap1 struct {
	RootRef execroot.RootRef
	LiabVar execvar.VarRec
}

type EvalSnap2 struct {
	RootRef    execroot.RootRef
	LiabMode   execvar.Mode
	StructVars []execvar.StructRec
	LinearVars []execvar.LinearRec
}

type EvalSnap3 struct {
	RootRef     execroot.RootRef
	StructVars  map[symbol.ADT]execvar.StructRec
	StructExprs map[symbol.ADT]typedef.ExprRec
	LinearVars  map[symbol.ADT]execvar.LinearRec
	LinearExprs map[symbol.ADT]typedef.ExprRec
}

type service struct {
	execEvalRepo   Repo
	execEvalBroker Broker
	execVarRepo    execvar.Repo
	commRepo       execcomm.Repo
	typeDefRepo    typedef.Repo
	procEvalRepo   procexeceval.Repo
	termDecRepo    termdec.Repo
	execRootRepo   execroot.Repo
	transactor     tx.Transactor
	log            *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

type ServiceDeps struct {
	ExecEvalRepo   Repo
	ExecEvalBroker Broker
	ExecVarRepo    execvar.Repo
	ExecCommRepo   execcomm.Repo
	TypeDefRepo    typedef.Repo
	ProcEvalRepo   procexeceval.Repo
	TermDecRepo    termdec.Repo
	ExecRootRepo   execroot.Repo
	Transactor     tx.Transactor
	Logger         *slog.Logger
}

func NewService(deps ServiceDeps) API {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{
		deps.ExecEvalRepo, deps.ExecEvalBroker, deps.ExecVarRepo,
		deps.ExecCommRepo, deps.TypeDefRepo, deps.ProcEvalRepo, deps.TermDecRepo,
		deps.ExecRootRepo,
		deps.Transactor, deps.Logger.With(name),
	}
}

func (s *service) Run(ctx context.Context, spec EvalSpec) (_ execroot.RootRef, err error) {
	specAttr := slog.Any("spec", spec)
	s.log.Debug("creation started", specAttr)
	var termDec termdec.DecRec
	getErr1 := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		termDec, err = s.termDecRepo.GetRecByQN(txCtx, conn, spec.TermQN)
		return err
	})
	if getErr1 != nil {
		s.log.Error("creation failed", specAttr)
		return execroot.RootRef{}, getErr1
	}
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars))
	for _, assetVar := range spec.AssetVars {
		if assetVar.TermQN == spec.LiabVar.TermQN {
			continue
		}
		assetQNs = append(assetQNs, assetVar.TermQN)
	}
	var assetEvals map[uniqsym.ADT]EvalSnap1
	getErr2 := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		assetEvals, err = s.execEvalRepo.GetSnapMapByQNs(txCtx, conn, assetQNs)
		return err
	})
	if getErr2 != nil {
		s.log.Error("creation failed", specAttr)
		return execroot.RootRef{}, getErr2
	}
	newEval := EvalRec{RootRef: execroot.New(), TermQN: spec.TermQN, LiabMode: execvar.StructMode}
	// Reserve the provider-side communication identity in its liability variable.
	// The Communication record itself is created when a client execution attaches.
	newCommID := identity.New()
	newLiabVar := execvar.StructRec{
		RootRef: newEval.RootRef,
		CommID:  newCommID,
		ChnlID:  identity.New(),
		ChnlPH:  spec.LiabVar.ChnlPH,
		ChnlBS:  execvar.LiabSide,
		ExprCK:  termDec.LiabVar.ExprCK,
	}
	newAssetVars := make([]execvar.VarRec, 0, len(spec.AssetVars)+1)
	// Materialize a Communication row for each attached asset endpoint.
	// A self-reference uses the new liability endpoint, so both endpoints and
	// its Communication are created atomically with this execution.
	newConns := make([]execcomm.CommRec, 0, len(assetEvals)+1)
	for _, assetVar := range spec.AssetVars {
		var commID identity.ADT
		var chnlID identity.ADT
		var exprCK contkey.ADT
		var providerExecRef execroot.RootRef
		assetExec, ok := assetEvals[assetVar.TermQN]
		if !ok && assetVar.TermQN == spec.LiabVar.TermQN {
			commID = newCommID
			chnlID = newLiabVar.ChnlID
			exprCK = newLiabVar.ExprCK
			providerExecRef = newEval.RootRef
		} else {
			commID = assetExec.LiabVar.GetCommID()
			providerExecRef = assetExec.RootRef
			chnlID = assetExec.LiabVar.GetChnlID()
			exprCK = assetExec.LiabVar.GetExprCK()
		}
		if ok || assetVar.TermQN == spec.LiabVar.TermQN {
			newConns = append(newConns, execcomm.CommRec{
				CommID:          commID,
				ProviderRootRef: providerExecRef,
				ClientRootRef:   newEval.RootRef,
				OffsetNr:        seqnum.Zero,
			})
		}
		newAssetVars = append(newAssetVars, execvar.StructRec{
			RootRef: newEval.RootRef,
			CommID:  commID,
			ChnlID:  chnlID,
			ChnlPH:  assetVar.ChnlPH,
			ChnlBS:  execvar.AssetSide,
			ExprCK:  exprCK,
		})
	}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		err = s.execEvalRepo.AddRec(txCtx, conn, newEval)
		if err != nil {
			return err
		}
		err = s.execVarRepo.AddRecs(txCtx, conn, append(newAssetVars, newLiabVar))
		if err != nil {
			return err
		}
		for _, newConn := range newConns {
			if err = s.commRepo.AddRec(txCtx, conn, newConn); err != nil {
				return err
			}
		}
		return nil
	})
	if transactErr != nil {
		s.log.Error("creation failed", specAttr)
		return execroot.RootRef{}, transactErr
	}
	s.log.Debug("creation succeed", slog.Any("ref", newEval.RootRef))
	return newEval.RootRef, nil
}

func (s *service) Spawn(ctx context.Context, spec execstep.StepSpec) (_ execroot.RootRef, err error) {
	refAttr := slog.Any("ref", spec.RootRef)
	s.log.Debug("proc spawning started", refAttr, slog.Any("expr", spec.PoolExpr))
	newEval := procexeceval.EvalRec{RootRef: execroot.New(), LiabMode: execvar.LinearMode}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		return s.procEvalRepo.AddRec(txCtx, conn, newEval)
	})
	if transactErr != nil {
		s.log.Error("proc spawning failed", refAttr)
		return execroot.RootRef{}, transactErr
	}
	s.log.Debug("proc spawning succeed", refAttr, slog.Any("proc", newEval.RootRef))
	return newEval.RootRef, nil
}

func (s *service) Take(ctx context.Context, spec execstep.StepSpec) (err error) {
	refAttr := slog.Any("ref", spec.RootRef)
	s.log.Debug("step taking started", refAttr, slog.Any("expr", spec.PoolExpr))
	evalSnap, retErr := s.retrieveSnap(ctx, spec.RootRef)
	if retErr != nil {
		s.log.Error("step taking failed", refAttr)
		return retErr
	}
	evalMod, evalEff, commMod, takeErr := s.takeStep(ctx, evalSnap, spec.PoolExpr)
	if takeErr != nil {
		s.log.Error("step taking failed", refAttr)
		return takeErr
	}
	// Every changed execution root is guarded by its observed root revision.
	// Communication changes always include the provider execution root, regardless of
	// which endpoint initiated the step.
	execRefs := map[string]execroot.RootRef{}
	if commMod.ProviderExecID.IsEmpty() {
		// Non-communication steps still mutate the executing Execution root.
		execRefs[evalSnap.RootRef.ExecID.String()] = evalSnap.RootRef
	} else {
		// A communication step always mutates the provider execution root. The
		// initiating client root is added only if its own variables change.
		providerExecRef := execroot.RootRef{
			ExecID: commMod.ProviderExecID,
			ExecRN: commMod.ProviderExecRN,
		}
		execRefs[providerExecRef.ExecID.String()] = providerExecRef
	}
	for _, variable := range evalMod.Vars {
		ref := variable.GetRootRef()
		// The Act's RootRef records the revision that produced the Act. For
		// provider-side writes, use the live provider revision captured with
		// the Communication snapshot instead of that historical Act revision.
		if !commMod.ProviderExecID.IsEmpty() && ref.ExecID == commMod.ProviderExecID {
			continue
		}
		if !ref.ExecID.IsEmpty() {
			execRefs[ref.ExecID.String()] = ref
		}
	}
	execKeys := make([]string, 0, len(execRefs))
	for key := range execRefs {
		execKeys = append(execKeys, key)
	}
	slices.Sort(execKeys)

	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		// CAS all affected aggregate roots before writing any child state.
		for _, key := range execKeys {
			if touchErr := s.execRootRepo.TouchRef(txCtx, conn, execRefs[key]); touchErr != nil {
				return touchErr
			}
		}
		if err = s.commRepo.ModifyRec(txCtx, conn, commMod); err != nil {
			return err
		}
		return s.execVarRepo.AddRecs(txCtx, conn, evalMod.Vars)
	})
	if transactErr != nil {
		s.log.Error("step taking failed", refAttr)
		return transactErr
	}
	for _, step := range evalEff.Steps {
		sendErr := s.execEvalBroker.SendSpec(step)
		if sendErr != nil {
			s.log.Error("step taking failed", refAttr)
			return sendErr
		}
	}
	return nil
}

func (s *service) takeStep(
	ctx context.Context,
	evalSnap EvalSnap3,
	expr ExprSpec,
) (
	evalMod EvalMod,
	evalEff EvalEff,
	commMod execcomm.CommMod,
	err error,
) {
	execAttr := slog.Any("rootRef", evalSnap.RootRef)
	switch termExpr := expr.(type) {
	case AcceptSpec:
		commChnl, ok := evalSnap.StructVars[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCfg(termExpr.CommChnlPH)
		}
		// вычисляем следующее состояние
		typeExpr, ok := evalSnap.StructExprs[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCtx(termExpr.CommChnlPH)
		}
		nextExprCK := typeExpr.(typedef.ProdRec).Next()
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
			return evalMod, evalEff, commMod, getErr
		}
		commMod.CommID = commSnap.CommID
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		commAttr := slog.Any("commID", commSnap.CommID)
		subscription := commSnap.NextAct()
		if subscription != nil {
			commMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступодателя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExpr: AcceptRec{
					ContChnlPH: commChnl.ChnlPH,
					ContExpr:   termExpr.ContExpr,
				},
			})
			s.log.Debug("taking half done", execAttr, commAttr)
			return evalMod, evalEff, commMod, nil
		}
		acquisition, ok := subscription.(commact.SubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение доступодателя
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: commChnl.RootRef,
			CommID:  commChnl.CommID,
			ChnlID:  newChnlID,
			ChnlPH:  commChnl.ChnlPH,
			ChnlBS:  commChnl.ChnlBS,
			ExprCK:  nextExprCK,
		})
		if termExpr.ContExpr != nil {
			// шедулим продолжение доступодателя
			evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
				RootRef:  evalSnap.RootRef,
				PoolExpr: termExpr.ContExpr,
			})
		}
		switch exprRec := acquisition.ContExpr.(type) {
		case AcquireRec:
			// сдвигаем офсет коммуникации
			// вяжем продолжение доступополучателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: acquisition.RootRef,
				CommID:  acquisition.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  exprRec.ContChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExprCK:  nextExprCK,
			})
			if exprRec.ContExpr != nil {
				// шедулим продолжение доступополучателя
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef:  acquisition.RootRef,
					PoolExpr: exprRec.ContExpr,
				})
			}
		default:
			panic(ErrRecTypeUnexpected(acquisition.ContExpr))
		}
		s.log.Debug("step taking succeed", execAttr, commAttr)
		return evalMod, evalEff, commMod, nil
	case AcquireSpec:
		commChnl, ok := evalSnap.StructVars[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCfg(termExpr.CommChnlPH)
		}
		// вычисляем следующее состояние
		typeExpr, ok := evalSnap.StructExprs[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCtx(termExpr.CommChnlPH)
		}
		nextExprCK := typeExpr.(typedef.ProdRec).Next()
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
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, getErr
		}
		commMod.CommID = commSnap.CommID
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		commAttr := slog.Any("commID", commSnap.CommID)
		subscription := commSnap.NextAct()
		if subscription != nil {
			commMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступополучателя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExpr: AcquireRec{
					ContChnlPH: commChnl.ChnlPH,
					ContExpr:   termExpr.ContExpr,
				},
			})
			s.log.Debug("taking half done", execAttr, commAttr)
			return evalMod, evalEff, commMod, nil
		}
		acception, ok := subscription.(commact.SubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение доступополучателя
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: commChnl.RootRef,
			CommID:  commChnl.CommID,
			ChnlID:  newChnlID,
			ChnlPH:  commChnl.ChnlPH,
			ChnlBS:  commChnl.ChnlBS,
			ExprCK:  nextExprCK,
		})
		if termExpr.ContExpr != nil {
			// шедулим продолжение доступополучателя
			evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
				RootRef:  evalSnap.RootRef,
				PoolExpr: termExpr.ContExpr,
			})
		}
		switch exprRec := acception.ContExpr.(type) {
		case AcceptRec:
			// сдвигаем офсет коммуникации
			// вяжем продолжение доступодателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: acception.RootRef,
				CommID:  acception.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  exprRec.ContChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExprCK:  nextExprCK,
			})
			if exprRec.ContExpr != nil {
				// шедулим продолжение доступодателя
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef:  acception.RootRef,
					PoolExpr: exprRec.ContExpr,
				})
			}
		default:
			panic(ErrRecTypeUnexpected(acception.ContExpr))
		}
		s.log.Debug("step taking succeed", execAttr, commAttr)
		return evalMod, evalEff, commMod, nil
	case ApplySpec:
		commChnl, ok := evalSnap.LinearVars[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCfg(termExpr.CommChnlPH)
		}
		// вычисляем следующее состояние
		typeExpr, ok := evalSnap.LinearExprs[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCtx(termExpr.CommChnlPH)
		}
		nextExprCK := typeExpr.(typedef.ProdRec).Next()
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
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, getErr
		}
		commMod.CommID = commSnap.CommID
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		commAttr := slog.Any("commID", commSnap.CommID)
		subscription := commSnap.NextAct()
		if subscription != nil {
			commMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку соискателя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExpr: ApplyRec{
					ContChnlPH: commChnl.ChnlPH,
					ProcTermQN: termExpr.ProcTermQN,
					ContExpr:   termExpr.ContExpr,
				},
			})
			s.log.Debug("taking half done", execAttr, commAttr)
			return evalMod, evalEff, commMod, nil
		}
		hiring, ok := subscription.(commact.SubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение соискателя
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: commChnl.RootRef,
			CommID:  commChnl.CommID,
			ChnlID:  newChnlID,
			ChnlPH:  commChnl.ChnlPH,
			ChnlBS:  commChnl.ChnlBS,
			ExprCK:  nextExprCK,
		})
		if termExpr.ContExpr != nil {
			// шедулим продолжение соискателя
			evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
				RootRef:  evalSnap.RootRef,
				PoolExpr: termExpr.ContExpr,
			})
		}
		switch exprRec := hiring.ContExpr.(type) {
		case HireRec:
			// вяжем продолжение нанимателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: hiring.RootRef,
				CommID:  hiring.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  exprRec.ContChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExprCK:  nextExprCK,
			})
			if exprRec.ContExpr != nil {
				// шедулим продолжение нанимателя
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef:  hiring.RootRef,
					PoolExpr: exprRec.ContExpr,
				})
			}
		default:
			panic(ErrRecTypeUnexpected(hiring.ContExpr))
		}
		s.log.Debug("step taking succeed", execAttr, commAttr)
		return evalMod, evalEff, commMod, nil
	case HireSpec:
		commChnl, ok := evalSnap.LinearVars[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCfg(termExpr.CommChnlPH)
		}
		// вычисляем следующее состояние
		typeExpr, ok := evalSnap.LinearExprs[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCtx(termExpr.CommChnlPH)
		}
		nextExprCK := typeExpr.(typedef.ProdRec).Next()
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
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, getErr
		}
		commMod.CommID = commSnap.CommID
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		commAttr := slog.Any("commID", commSnap.CommID)
		subscription := commSnap.NextAct()
		if subscription != nil {
			commMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку нанимателя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExpr: HireRec{
					ContChnlPH: commChnl.ChnlPH,
					ProcTermQN: termExpr.ProcTermQN,
					ContExpr:   termExpr.ContExpr,
				},
			})
			s.log.Debug("taking half done", execAttr, commAttr)
			return evalMod, evalEff, commMod, nil
		}
		application, ok := subscription.(commact.SubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение нанимателя
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: commChnl.RootRef,
			CommID:  commChnl.CommID,
			ChnlID:  newChnlID,
			ChnlPH:  commChnl.ChnlPH,
			ChnlBS:  commChnl.ChnlBS,
			ExprCK:  nextExprCK,
		})
		if termExpr.ContExpr != nil {
			// шедулим продолжение нанимателя
			evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
				RootRef:  evalSnap.RootRef,
				PoolExpr: termExpr.ContExpr,
			})
		}
		switch exprRec := application.ContExpr.(type) {
		case ApplyRec:
			// вяжем продолжение соискателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: application.RootRef,
				CommID:  application.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  exprRec.ContChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExprCK:  nextExprCK,
			})
			if exprRec.ContExpr != nil {
				// запускаем продолжение соискателя
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef:  application.RootRef,
					PoolExpr: exprRec.ContExpr,
				})
			}
		default:
			panic(ErrRecTypeUnexpected(application.ContExpr))
		}
		s.log.Debug("step taking succeed", execAttr, commAttr)
		return evalMod, evalEff, commMod, nil
	case ReleaseSpec:
		commChnl, ok := evalSnap.LinearVars[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCfg(termExpr.CommChnlPH)
		}
		// вычисляем следующее состояние
		typeExpr, ok := evalSnap.LinearExprs[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCtx(termExpr.CommChnlPH)
		}
		nextExprCK := typeExpr.(typedef.ProdRec).Next()
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
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, getErr
		}
		commMod.CommID = commSnap.CommID
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		commAttr := slog.Any("commID", commSnap.CommID)
		subscription := commSnap.NextAct()
		if subscription != nil {
			commMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступовозвращателя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExpr: ReleaseRec{
					ContChnlPH: commChnl.ChnlPH,
				},
			})
			s.log.Debug("taking half done", execAttr, commAttr)
			return evalMod, evalEff, commMod, nil
		}
		detaching, ok := subscription.(commact.SubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение доступовозвращателя
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: commChnl.RootRef,
			CommID:  commChnl.CommID,
			ChnlID:  newChnlID,
			ChnlPH:  commChnl.ChnlPH,
			ChnlBS:  commChnl.ChnlBS,
			ExprCK:  nextExprCK,
		})
		switch exprRec := detaching.ContExpr.(type) {
		case DetachRec:
			// вяжем продолжение доступопринимателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: detaching.RootRef,
				CommID:  detaching.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  exprRec.ContChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExprCK:  nextExprCK,
			})
		default:
			panic(ErrRecTypeUnexpected(detaching.ContExpr))
		}
		s.log.Debug("step taking succeed", execAttr, commAttr)
		return evalMod, evalEff, commMod, nil
	case DetachSpec:
		commChnl, ok := evalSnap.LinearVars[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCfg(termExpr.CommChnlPH)
		}
		// вычисляем следующее состояние
		typeExpr, ok := evalSnap.LinearExprs[termExpr.CommChnlPH]
		if !ok {
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, proctermdef.ErrMissingInCtx(termExpr.CommChnlPH)
		}
		nextExprCK := typeExpr.(typedef.ProdRec).Next()
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
			s.log.Error("step taking failed", execAttr)
			return evalMod, evalEff, commMod, getErr
		}
		commMod.CommID = commSnap.CommID
		commMod.ProviderExecID = commSnap.ProviderExecID
		commMod.ProviderExecRN = commSnap.ProviderExecRN
		commAttr := slog.Any("commID", commSnap.CommID)
		subscription := commSnap.NextAct()
		if subscription != nil {
			commMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступопринимателя
			commMod.Acts = append(commMod.Acts, commact.SubRec{
				CommID:  commSnap.CommID,
				RootRef: evalSnap.RootRef,
				ChnlID:  commChnl.ChnlID,
				ContExpr: DetachRec{
					ContChnlPH: commChnl.ChnlPH,
				},
			})
			s.log.Debug("taking half done", execAttr, commAttr)
			return evalMod, evalEff, commMod, nil
		}
		releasing, ok := subscription.(commact.SubRec)
		if !ok {
			panic(commact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: commChnl.RootRef,
			CommID:  commChnl.CommID,
			ChnlID:  newChnlID,
			ChnlPH:  commChnl.ChnlPH,
			ChnlBS:  commChnl.ChnlBS,
			ExprCK:  nextExprCK,
		})
		switch exprRec := releasing.ContExpr.(type) {
		case ReleaseRec:
			// вяжем продолжение доступовозвращателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: releasing.RootRef,
				CommID:  releasing.CommID,
				ChnlID:  newChnlID,
				ChnlPH:  exprRec.ContChnlPH,
				ChnlBS:  commChnl.ChnlBS,
				ExprCK:  nextExprCK,
			})
		default:
			panic(ErrRecTypeUnexpected(releasing.ContExpr))
		}
		s.log.Debug("step taking succeed", execAttr, commAttr)
		return evalMod, evalEff, commMod, nil
	default:
		panic(ErrSpecTypeUnexpected(expr))
	}
}

func (s *service) retrieveSnap(ctx context.Context, ref execroot.RootRef) (_ EvalSnap3, err error) {
	var evalSnap EvalSnap2
	getErr1 := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		evalSnap, err = s.execEvalRepo.GetSnapByRef(txCtx, conn, ref)
		return err
	})
	if getErr1 != nil {
		return EvalSnap3{}, getErr1
	}
	var structExprs map[symbol.ADT]typedef.ExprRec
	var linearExprs map[symbol.ADT]typedef.ExprRec
	getErr2 := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		structExprs, err = s.typeDefRepo.GetTypeMap(txCtx, conn, ExtractExprCKs(evalSnap.StructVars))
		if err != nil {
			return err
		}
		linearExprs, err = s.typeDefRepo.GetTypeMap(txCtx, conn, ExtractExprCKs(evalSnap.LinearVars))
		return err
	})
	if getErr2 != nil {
		return EvalSnap3{}, getErr2
	}
	return EvalSnap3{
		RootRef:     evalSnap.RootRef,
		StructVars:  execvar.ConvertRecsToRecMap(evalSnap.StructVars),
		StructExprs: structExprs,
		LinearVars:  execvar.ConvertRecsToRecMap(evalSnap.LinearVars),
		LinearExprs: linearExprs,
	}, nil
}
