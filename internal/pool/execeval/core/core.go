package core

import (
	"context"
	"log/slog"
	"reflect"
	"slices"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	connact "orglang/go-engine/internal/pool/connact/core"
	execconn "orglang/go-engine/internal/pool/execconn/core"
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
	RootRef    execroot.RootRef
	StructVars map[symbol.ADT]execvar.StructRec
	StructExps map[symbol.ADT]typedef.ExpRec
	LinearVars map[symbol.ADT]execvar.LinearRec
	LinearExps map[symbol.ADT]typedef.ExpRec
}

type service struct {
	execEvalRepo   Repo
	execEvalBroker Broker
	execVarRepo    execvar.Repo
	connRepo       execconn.Repo
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
	ExecConnRepo   execconn.Repo
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
		deps.ExecConnRepo, deps.TypeDefRepo, deps.ProcEvalRepo, deps.TermDecRepo,
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
	// Reserve the provider-side connection identity in its liability variable.
	// The Connection record itself is created when a client Computation attaches.
	newConnID := identity.New()
	newLiabVar := execvar.StructRec{
		RootRef: newEval.RootRef,
		ConnID:  newConnID,
		ChnlID:  identity.New(),
		ChnlPH:  spec.LiabVar.ChnlPH,
		ChnlBS:  execvar.LiabSide,
		ExpCK:   termDec.LiabVar.ExpCK,
	}
	newAssetVars := make([]execvar.VarRec, 0, len(spec.AssetVars)+1)
	// Materialize a Connection row for each attached asset endpoint.
	// A self-reference uses the new liability endpoint, so both endpoints and
	// its Connection are created atomically with this Computation.
	newConns := make([]execconn.ConnRec, 0, len(assetEvals)+1)
	for _, assetVar := range spec.AssetVars {
		var connID identity.ADT
		var chnlID identity.ADT
		var expCK contkey.ADT
		var providerRootRef execroot.RootRef
		assetExec, ok := assetEvals[assetVar.TermQN]
		if !ok && assetVar.TermQN == spec.LiabVar.TermQN {
			connID = newConnID
			chnlID = newLiabVar.ChnlID
			expCK = newLiabVar.ExpCK
			providerRootRef = newEval.RootRef
		} else {
			connID = assetExec.LiabVar.GetConnID()
			providerRootRef = assetExec.RootRef
			chnlID = assetExec.LiabVar.GetChnlID()
			expCK = assetExec.LiabVar.GetExpCK()
		}
		if ok || assetVar.TermQN == spec.LiabVar.TermQN {
			newConns = append(newConns, execconn.ConnRec{
				ConnID:          connID,
				ProviderRootRef: providerRootRef,
				ClientRootRef:   newEval.RootRef,
				OffsetNr:        seqnum.Zero,
			})
		}
		newAssetVars = append(newAssetVars, execvar.StructRec{
			RootRef: newEval.RootRef,
			ConnID:  connID,
			ChnlID:  chnlID,
			ChnlPH:  assetVar.ChnlPH,
			ChnlBS:  execvar.AssetSide,
			ExpCK:   expCK,
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
			if err = s.connRepo.AddRec(txCtx, conn, newConn); err != nil {
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
	s.log.Debug("proc spawning started", refAttr, slog.Any("exp", spec.PoolExp))
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
	s.log.Debug("step taking started", refAttr, slog.Any("exp", spec.PoolExp))
	evalSnap, retErr := s.retrieveSnap(ctx, spec.RootRef)
	if retErr != nil {
		s.log.Error("step taking failed", refAttr)
		return retErr
	}
	evalMod, evalEff, connMod, takeErr := s.takeStep(ctx, evalSnap, spec.PoolExp)
	if takeErr != nil {
		s.log.Error("step taking failed", refAttr)
		return takeErr
	}
	// Every changed Computation aggregate is guarded by its observed root revision.
	// Connection changes always include the provider aggregate, regardless of
	// which endpoint initiated the step.
	rootRefs := map[string]execroot.RootRef{}
	if connMod.ProviderExecID.IsEmpty() {
		// Non-connection steps still mutate the executing Computation root.
		rootRefs[evalSnap.RootRef.RootID.String()] = evalSnap.RootRef
	} else {
		// A connection step always mutates the provider aggregate. The
		// initiating client root is added only if its own variables change.
		providerRef := execroot.RootRef{
			RootID: connMod.ProviderExecID,
			RootRN: connMod.ProviderExecRN,
		}
		rootRefs[providerRef.RootID.String()] = providerRef
	}
	for _, variable := range evalMod.Vars {
		ref := variable.GetRootRef()
		// The Act's RootRef records the revision that produced the Act. For
		// provider-side writes, use the live provider revision captured with
		// the Connection snapshot instead of that historical Act revision.
		if !connMod.ProviderExecID.IsEmpty() && ref.RootID == connMod.ProviderExecID {
			continue
		}
		if !ref.RootID.IsEmpty() {
			rootRefs[ref.RootID.String()] = ref
		}
	}
	compKeys := make([]string, 0, len(rootRefs))
	for key := range rootRefs {
		compKeys = append(compKeys, key)
	}
	slices.Sort(compKeys)

	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		// CAS all affected aggregate roots before writing any child state.
		for _, key := range compKeys {
			if touchErr := s.execRootRepo.TouchRef(txCtx, conn, rootRefs[key]); touchErr != nil {
				return touchErr
			}
		}
		if err = s.connRepo.ModifyRec(txCtx, conn, connMod); err != nil {
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
	exp ExpSpec,
) (
	evalMod EvalMod,
	evalEff EvalEff,
	connMod execconn.ConnMod,
	err error,
) {
	rootAttr := slog.Any("rootRef", evalSnap.RootRef)
	switch termExp := exp.(type) {
	case AcceptSpec:
		connChnl, ok := evalSnap.StructVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := evalSnap.StructExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typedef.ProdRec).Next()
		// получаем снепшот соединения
		var connSnap execconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, execconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderExecID = connSnap.ProviderExecID
		connMod.ProviderExecRN = connSnap.ProviderExecRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступодателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				RootRef: evalSnap.RootRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: AcceptRec{
					ContChnlPH: connChnl.ChnlPH,
					ContExp:    termExp.ContExp,
				},
			})
			s.log.Debug("taking half done", rootAttr, connAttr)
			return evalMod, evalEff, connMod, nil
		}
		acquisition, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение доступодателя
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: connChnl.RootRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		if termExp.ContExp != nil {
			// шедулим продолжение доступодателя
			evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
				RootRef: evalSnap.RootRef,
				PoolExp: termExp.ContExp,
			})
		}
		switch expRec := acquisition.ContExp.(type) {
		case AcquireRec:
			// сдвигаем офсет коммуникации
			// вяжем продолжение доступополучателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: acquisition.RootRef,
				ConnID:  acquisition.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение доступополучателя
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: acquisition.RootRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(ErrRecTypeUnexpected(acquisition.ContExp))
		}
		s.log.Debug("step taking succeed", rootAttr, connAttr)
		return evalMod, evalEff, connMod, nil
	case AcquireSpec:
		connChnl, ok := evalSnap.StructVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := evalSnap.StructExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typedef.ProdRec).Next()
		// получаем снепшот коммуникации
		var connSnap execconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, execconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderExecID = connSnap.ProviderExecID
		connMod.ProviderExecRN = connSnap.ProviderExecRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступополучателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				RootRef: evalSnap.RootRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: AcquireRec{
					ContChnlPH: connChnl.ChnlPH,
					ContExp:    termExp.ContExp,
				},
			})
			s.log.Debug("taking half done", rootAttr, connAttr)
			return evalMod, evalEff, connMod, nil
		}
		acception, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение доступополучателя
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: connChnl.RootRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		if termExp.ContExp != nil {
			// шедулим продолжение доступополучателя
			evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
				RootRef: evalSnap.RootRef,
				PoolExp: termExp.ContExp,
			})
		}
		switch expRec := acception.ContExp.(type) {
		case AcceptRec:
			// сдвигаем офсет коммуникации
			// вяжем продолжение доступодателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: acception.RootRef,
				ConnID:  acception.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение доступодателя
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: acception.RootRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(ErrRecTypeUnexpected(acception.ContExp))
		}
		s.log.Debug("step taking succeed", rootAttr, connAttr)
		return evalMod, evalEff, connMod, nil
	case ApplySpec:
		connChnl, ok := evalSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := evalSnap.LinearExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typedef.ProdRec).Next()
		// получаем снепшот коммуникации
		var connSnap execconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, execconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderExecID = connSnap.ProviderExecID
		connMod.ProviderExecRN = connSnap.ProviderExecRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку соискателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				RootRef: evalSnap.RootRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: ApplyRec{
					ContChnlPH: connChnl.ChnlPH,
					ProcTermQN: termExp.ProcTermQN,
					ContExp:    termExp.ContExp,
				},
			})
			s.log.Debug("taking half done", rootAttr, connAttr)
			return evalMod, evalEff, connMod, nil
		}
		hiring, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение соискателя
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: connChnl.RootRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		if termExp.ContExp != nil {
			// шедулим продолжение соискателя
			evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
				RootRef: evalSnap.RootRef,
				PoolExp: termExp.ContExp,
			})
		}
		switch expRec := hiring.ContExp.(type) {
		case HireRec:
			// вяжем продолжение нанимателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: hiring.RootRef,
				ConnID:  hiring.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение нанимателя
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: hiring.RootRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(ErrRecTypeUnexpected(hiring.ContExp))
		}
		s.log.Debug("step taking succeed", rootAttr, connAttr)
		return evalMod, evalEff, connMod, nil
	case HireSpec:
		connChnl, ok := evalSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := evalSnap.LinearExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typedef.ProdRec).Next()
		// получаем снепшот коммуникации
		var connSnap execconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, execconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderExecID = connSnap.ProviderExecID
		connMod.ProviderExecRN = connSnap.ProviderExecRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку нанимателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				RootRef: evalSnap.RootRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: HireRec{
					ContChnlPH: connChnl.ChnlPH,
					ProcTermQN: termExp.ProcTermQN,
					ContExp:    termExp.ContExp,
				},
			})
			s.log.Debug("taking half done", rootAttr, connAttr)
			return evalMod, evalEff, connMod, nil
		}
		application, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение нанимателя
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: connChnl.RootRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		if termExp.ContExp != nil {
			// шедулим продолжение нанимателя
			evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
				RootRef: evalSnap.RootRef,
				PoolExp: termExp.ContExp,
			})
		}
		switch expRec := application.ContExp.(type) {
		case ApplyRec:
			// вяжем продолжение соискателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: application.RootRef,
				ConnID:  application.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if expRec.ContExp != nil {
				// запускаем продолжение соискателя
				evalEff.Steps = append(evalEff.Steps, execstep.StepSpec{
					RootRef: application.RootRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(ErrRecTypeUnexpected(application.ContExp))
		}
		s.log.Debug("step taking succeed", rootAttr, connAttr)
		return evalMod, evalEff, connMod, nil
	case ReleaseSpec:
		connChnl, ok := evalSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := evalSnap.LinearExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typedef.ProdRec).Next()
		// получаем снепшот коммуникации
		var connSnap execconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, execconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderExecID = connSnap.ProviderExecID
		connMod.ProviderExecRN = connSnap.ProviderExecRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступовозвращателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				RootRef: evalSnap.RootRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: ReleaseRec{
					ContChnlPH: connChnl.ChnlPH,
				},
			})
			s.log.Debug("taking half done", rootAttr, connAttr)
			return evalMod, evalEff, connMod, nil
		}
		detaching, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение доступовозвращателя
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: connChnl.RootRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		switch expRec := detaching.ContExp.(type) {
		case DetachRec:
			// вяжем продолжение доступопринимателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: detaching.RootRef,
				ConnID:  detaching.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
		default:
			panic(ErrRecTypeUnexpected(detaching.ContExp))
		}
		s.log.Debug("step taking succeed", rootAttr, connAttr)
		return evalMod, evalEff, connMod, nil
	case DetachSpec:
		connChnl, ok := evalSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := evalSnap.LinearExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typedef.ProdRec).Next()
		// получаем снепшот коммуникации
		var connSnap execconn.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, execconn.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", rootAttr)
			return evalMod, evalEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderExecID = connSnap.ProviderExecID
		connMod.ProviderExecRN = connSnap.ProviderExecRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступопринимателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				RootRef: evalSnap.RootRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: DetachRec{
					ContChnlPH: connChnl.ChnlPH,
				},
			})
			s.log.Debug("taking half done", rootAttr, connAttr)
			return evalMod, evalEff, connMod, nil
		}
		releasing, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
			RootRef: connChnl.RootRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		switch expRec := releasing.ContExp.(type) {
		case ReleaseRec:
			// вяжем продолжение доступовозвращателя
			evalMod.Vars = append(evalMod.Vars, execvar.LinearRec{
				RootRef: releasing.RootRef,
				ConnID:  releasing.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
		default:
			panic(ErrRecTypeUnexpected(releasing.ContExp))
		}
		s.log.Debug("step taking succeed", rootAttr, connAttr)
		return evalMod, evalEff, connMod, nil
	default:
		panic(ErrSpecTypeUnexpected(exp))
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
	var structExps map[symbol.ADT]typedef.ExpRec
	var linearExps map[symbol.ADT]typedef.ExpRec
	getErr2 := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		structExps, err = s.typeDefRepo.GetTypeMap(txCtx, conn, ExtractExpCKs(evalSnap.StructVars))
		if err != nil {
			return err
		}
		linearExps, err = s.typeDefRepo.GetTypeMap(txCtx, conn, ExtractExpCKs(evalSnap.LinearVars))
		return err
	})
	if getErr2 != nil {
		return EvalSnap3{}, getErr2
	}
	return EvalSnap3{
		RootRef:    evalSnap.RootRef,
		StructVars: execvar.ConvertRecsToRecMap(evalSnap.StructVars),
		StructExps: structExps,
		LinearVars: execvar.ConvertRecsToRecMap(evalSnap.LinearVars),
		LinearExps: linearExps,
	}, nil
}
