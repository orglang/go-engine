package core

import (
	"context"
	"log/slog"
	"reflect"
	"slices"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/implsem"
	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	compconncore "orglang/go-engine/internal/pool/compconn/core"
	connact "orglang/go-engine/internal/pool/connact/core"
	compstep "orglang/go-engine/internal/pool/compstep/core"
	compvar "orglang/go-engine/internal/pool/compvar/core"
	termdef "orglang/go-engine/internal/pool/termdef/core"
	termexp "orglang/go-engine/internal/pool/termexp/core"
	typeexp "orglang/go-engine/internal/pool/typeexp/core"

	proccompexec "orglang/go-engine/internal/proc/compexec/core"
	proctermdef "orglang/go-engine/internal/proc/termdef/core"
)

type API interface {
	Run(ExecSpec) (compsem.SemRef, error) // aka Create
	Take(compstep.StepSpec) error
	Spawn(compstep.StepSpec) (compsem.SemRef, error)
}

type Broker interface {
	Subscribe(api API)
	SendSpec(compstep.StepSpec) error
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
	connRepo       compconncore.Repo
	typeExpRepo    typeexp.Repo
	procExecRepo   proccompexec.Repo
	termDefRepo    termdef.Repo
	implSemRepo    implsem.Repo
	compSemRepo    compsem.Repo
	transactor     tx.Transactor
	log            *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

type ServiceDeps struct {
	CompExecRepo   Repo
	CompExecBroker Broker
	CompVarRepo    compvar.Repo
	ConnRepo   compconncore.Repo
	TypeExpRepo    typeexp.Repo
	ProcExecRepo   proccompexec.Repo
	TermDefRepo    termdef.Repo
	ImplSemRepo    implsem.Repo
	CompSemRepo    compsem.Repo
	Transactor     tx.Transactor
	Logger         *slog.Logger
}

func NewService(deps ServiceDeps) API {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{
		deps.CompExecRepo, deps.CompExecBroker, deps.CompVarRepo,
		deps.ConnRepo, deps.TypeExpRepo, deps.ProcExecRepo, deps.TermDefRepo,
		deps.ImplSemRepo, deps.CompSemRepo,
		deps.Transactor, deps.Logger.With(name),
	}
}

func (s *service) Run(spec ExecSpec) (_ compsem.SemRef, err error) {
	ctx := context.Background()
	specAttr := slog.Any("spec", spec)
	s.log.Debug("creation started", specAttr)
	var termDec termdef.DefRec
	getErr1 := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		termDec, err = s.termDefRepo.GetRecByQN(txCtx, conn, spec.TermQN)
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
	getErr2 := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		assetExecs, err = s.compExecRepo.GetSnapMapByQNs(txCtx, conn, assetQNs)
		return err
	})
	if getErr2 != nil {
		s.log.Error("creation failed", specAttr)
		return compsem.SemRef{}, getErr2
	}
	newExec := ExecRec{CompRef: compsem.New(), LiabMode: compvar.StructMode}
	// Reserve the provider-side connection identity in its liability variable.
	// The Connection record itself is created when a client Computation attaches.
	newConnID := identity.New()
	newImpl := implsem.SemRec{ImplQN: spec.TermQN, ImplID: newExec.CompRef.CompID}
	newLiabVar := compvar.StructRec{
		CompRef: newExec.CompRef,
		ConnID:  newConnID,
		ChnlID:  identity.New(),
		ChnlPH:  spec.LiabVar.ChnlPH,
		ChnlBS:  compvar.LiabSide,
		ExpCK:   termDec.LiabVar.ExpCK,
	}
	newAssetVars := make([]compvar.VarRec, 0, len(spec.AssetVars)+1)
	// Materialize a Connection row for each attached asset endpoint.
	// A self-reference uses the new liability endpoint, so both endpoints and
	// its Connection are created atomically with this Computation.
	newConns := make([]compconncore.ConnRec, 0, len(assetExecs)+1)
	for _, assetVar := range spec.AssetVars {
		var connID identity.ADT
		var chnlID identity.ADT
		var expCK contkey.ADT
		var providerCompRef compsem.SemRef
		assetExec, ok := assetExecs[assetVar.TermQN]
		if !ok && assetVar.TermQN == spec.LiabVar.TermQN {
			connID = newConnID
			chnlID = newLiabVar.ChnlID
			expCK = newLiabVar.ExpCK
			providerCompRef = newExec.CompRef
		} else {
			connID = assetExec.LiabVar.GetConnID()
			providerCompRef = assetExec.CompRef
			chnlID = assetExec.LiabVar.GetChnlID()
			expCK = assetExec.LiabVar.GetExpCK()
		}
		if ok || assetVar.TermQN == spec.LiabVar.TermQN {
			newConns = append(newConns, compconncore.ConnRec{
				ConnID:          connID,
				ProviderCompRef: providerCompRef,
				ClientCompRef:   newExec.CompRef,
				OffsetNr:        seqnum.Zero,
			})
		}
		newAssetVars = append(newAssetVars, compvar.StructRec{
			CompRef: newExec.CompRef,
			ConnID:  connID,
			ChnlID:  chnlID,
			ChnlPH:  assetVar.ChnlPH,
			ChnlBS:  compvar.AssetSide,
			ExpCK:   expCK,
		})
	}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		err = s.implSemRepo.AddRec(txCtx, conn, newImpl)
		if err != nil {
			return err
		}
		err = s.compExecRepo.AddRec(txCtx, conn, newExec)
		if err != nil {
			return err
		}
		err = s.compVarRepo.AddRecs(txCtx, conn, append(newAssetVars, newLiabVar))
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
		return compsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", slog.Any("ref", newExec.CompRef))
	return newExec.CompRef, nil
}

func (s *service) Spawn(spec compstep.StepSpec) (_ compsem.SemRef, err error) {
	ctx := context.Background()
	refAttr := slog.Any("ref", spec.CompRef)
	s.log.Debug("proc spawning started", refAttr, slog.Any("exp", spec.PoolExp))
	newExec := proccompexec.ExecRec{CompRef: compsem.New(), LiabMode: compvar.LinearMode}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		return s.procExecRepo.AddRec(txCtx, conn, newExec)
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
	execMod, execEff, connMod, takeErr := s.takeStep(execSnap, spec.PoolExp)
	if takeErr != nil {
		s.log.Error("step taking failed", refAttr)
		return takeErr
	}
	// Every changed Computation aggregate is guarded by its observed root revision.
	// Connection changes always include the provider aggregate, regardless of
	// which endpoint initiated the step.
	compRefs := map[string]compsem.SemRef{}
	if connMod.ProviderCompID.IsEmpty() {
		// Non-connection steps still mutate the executing Computation root.
		compRefs[execSnap.CompRef.CompID.String()] = execSnap.CompRef
	} else {
		// A connection step always mutates the provider aggregate. The
		// initiating client root is added only if its own variables change.
		providerRef := compsem.SemRef{
			CompID: connMod.ProviderCompID,
			CompRN: connMod.ProviderCompRN,
		}
		compRefs[providerRef.CompID.String()] = providerRef
	}
	for _, variable := range execMod.Vars {
		ref := variable.GetCompRef()
		// The Act's CompRef records the revision that produced the Act. For
		// provider-side writes, use the live provider revision captured with
		// the Connection snapshot instead of that historical Act revision.
		if !connMod.ProviderCompID.IsEmpty() && ref.CompID == connMod.ProviderCompID {
			continue
		}
		if !ref.CompID.IsEmpty() {
			compRefs[ref.CompID.String()] = ref
		}
	}
	compKeys := make([]string, 0, len(compRefs))
	for key := range compRefs {
		compKeys = append(compKeys, key)
	}
	slices.Sort(compKeys)

	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		// CAS all affected aggregate roots before writing any child state.
		for _, key := range compKeys {
			if touchErr := s.compSemRepo.TouchRef(txCtx, conn, compRefs[key]); touchErr != nil {
				return touchErr
			}
		}
		if err = s.connRepo.ModifyRec(txCtx, conn, connMod); err != nil {
			return err
		}
		return s.compVarRepo.AddRecs(txCtx, conn, execMod.Vars)
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

func (s *service) takeStep(
	execSnap ExecSnap3,
	exp termexp.ExpSpec,
) (
	execMod ExecMod,
	execEff ExecEff,
	connMod compconncore.ConnMod,
	err error,
) {
	ctx := context.Background()
	compAttr := slog.Any("compRef", execSnap.CompRef)
	switch termExp := exp.(type) {
	case termexp.AcceptSpec:
		connChnl, ok := execSnap.StructVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := execSnap.StructExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typeexp.ProdRec).Next()
		// получаем снепшот соединения
		var connSnap compconncore.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconncore.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступодателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: termexp.AcceptRec{
					ContChnlPH: connChnl.ChnlPH,
					ContExp:    termExp.ContExp,
				},
			})
			s.log.Debug("taking half done", compAttr, connAttr)
			return execMod, execEff, connMod, nil
		}
		acquisition, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение доступодателя
		execMod.Vars = append(execMod.Vars, compvar.LinearRec{
			CompRef: connChnl.CompRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		if termExp.ContExp != nil {
			// шедулим продолжение доступодателя
			execEff.Steps = append(execEff.Steps, compstep.StepSpec{
				CompRef: execSnap.CompRef,
				PoolExp: termExp.ContExp,
			})
		}
		switch expRec := acquisition.ContExp.(type) {
		case termexp.AcquireRec:
			// сдвигаем офсет коммуникации
			// вяжем продолжение доступополучателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				CompRef: acquisition.CompRef,
				ConnID:  acquisition.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение доступополучателя
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: acquisition.CompRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(termexp.ErrRecTypeUnexpected(acquisition.ContExp))
		}
		s.log.Debug("step taking succeed", compAttr, connAttr)
		return execMod, execEff, connMod, nil
	case termexp.AcquireSpec:
		connChnl, ok := execSnap.StructVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := execSnap.StructExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var connSnap compconncore.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconncore.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступополучателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: termexp.AcquireRec{
					ContChnlPH: connChnl.ChnlPH,
					ContExp:    termExp.ContExp,
				},
			})
			s.log.Debug("taking half done", compAttr, connAttr)
			return execMod, execEff, connMod, nil
		}
		acception, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение доступополучателя
		execMod.Vars = append(execMod.Vars, compvar.LinearRec{
			CompRef: connChnl.CompRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		if termExp.ContExp != nil {
			// шедулим продолжение доступополучателя
			execEff.Steps = append(execEff.Steps, compstep.StepSpec{
				CompRef: execSnap.CompRef,
				PoolExp: termExp.ContExp,
			})
		}
		switch expRec := acception.ContExp.(type) {
		case termexp.AcceptRec:
			// сдвигаем офсет коммуникации
			// вяжем продолжение доступодателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				CompRef: acception.CompRef,
				ConnID:  acception.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение доступодателя
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: acception.CompRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(termexp.ErrRecTypeUnexpected(acception.ContExp))
		}
		s.log.Debug("step taking succeed", compAttr, connAttr)
		return execMod, execEff, connMod, nil
	case termexp.ApplySpec:
		connChnl, ok := execSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := execSnap.LinearExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var connSnap compconncore.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconncore.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку соискателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: termexp.ApplyRec{
					ContChnlPH: connChnl.ChnlPH,
					ProcTermQN: termExp.ProcTermQN,
					ContExp:    termExp.ContExp,
				},
			})
			s.log.Debug("taking half done", compAttr, connAttr)
			return execMod, execEff, connMod, nil
		}
		hiring, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение соискателя
		execMod.Vars = append(execMod.Vars, compvar.LinearRec{
			CompRef: connChnl.CompRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		if termExp.ContExp != nil {
			// шедулим продолжение соискателя
			execEff.Steps = append(execEff.Steps, compstep.StepSpec{
				CompRef: execSnap.CompRef,
				PoolExp: termExp.ContExp,
			})
		}
		switch expRec := hiring.ContExp.(type) {
		case termexp.HireRec:
			// вяжем продолжение нанимателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				CompRef: hiring.CompRef,
				ConnID:  hiring.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if expRec.ContExp != nil {
				// шедулим продолжение нанимателя
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: hiring.CompRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(termexp.ErrRecTypeUnexpected(hiring.ContExp))
		}
		s.log.Debug("step taking succeed", compAttr, connAttr)
		return execMod, execEff, connMod, nil
	case termexp.HireSpec:
		connChnl, ok := execSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := execSnap.LinearExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var connSnap compconncore.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconncore.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку нанимателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: termexp.HireRec{
					ContChnlPH: connChnl.ChnlPH,
					ProcTermQN: termExp.ProcTermQN,
					ContExp:    termExp.ContExp,
				},
			})
			s.log.Debug("taking half done", compAttr, connAttr)
			return execMod, execEff, connMod, nil
		}
		application, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение нанимателя
		execMod.Vars = append(execMod.Vars, compvar.LinearRec{
			CompRef: connChnl.CompRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		if termExp.ContExp != nil {
			// шедулим продолжение нанимателя
			execEff.Steps = append(execEff.Steps, compstep.StepSpec{
				CompRef: execSnap.CompRef,
				PoolExp: termExp.ContExp,
			})
		}
		switch expRec := application.ContExp.(type) {
		case termexp.ApplyRec:
			// вяжем продолжение соискателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				CompRef: application.CompRef,
				ConnID:  application.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
			if expRec.ContExp != nil {
				// запускаем продолжение соискателя
				execEff.Steps = append(execEff.Steps, compstep.StepSpec{
					CompRef: application.CompRef,
					PoolExp: expRec.ContExp,
				})
			}
		default:
			panic(termexp.ErrRecTypeUnexpected(application.ContExp))
		}
		s.log.Debug("step taking succeed", compAttr, connAttr)
		return execMod, execEff, connMod, nil
	case termexp.ReleaseSpec:
		connChnl, ok := execSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := execSnap.LinearExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var connSnap compconncore.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconncore.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступовозвращателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: termexp.ReleaseRec{
					ContChnlPH: connChnl.ChnlPH,
				},
			})
			s.log.Debug("taking half done", compAttr, connAttr)
			return execMod, execEff, connMod, nil
		}
		detaching, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		// вяжем продолжение доступовозвращателя
		execMod.Vars = append(execMod.Vars, compvar.LinearRec{
			CompRef: connChnl.CompRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		switch expRec := detaching.ContExp.(type) {
		case termexp.DetachRec:
			// вяжем продолжение доступопринимателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				CompRef: detaching.CompRef,
				ConnID:  detaching.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
		default:
			panic(termexp.ErrRecTypeUnexpected(detaching.ContExp))
		}
		s.log.Debug("step taking succeed", compAttr, connAttr)
		return execMod, execEff, connMod, nil
	case termexp.DetachSpec:
		connChnl, ok := execSnap.LinearVars[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCfg(termExp.ConnChnlPH)
		}
		// вычисляем следующее состояние
		typeExp, ok := execSnap.LinearExps[termExp.ConnChnlPH]
		if !ok {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, proctermdef.ErrMissingInCtx(termExp.ConnChnlPH)
		}
		nextExpCK := typeExp.(typeexp.ProdRec).Next()
		// получаем снепшот коммуникации
		var connSnap compconncore.ConnSnap
		getErr := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
			connSnap, err = s.connRepo.GetSnapByQry(txCtx, conn, compconncore.ConnQry{
				ConnID: connChnl.ConnID,
				ChnlID: option.Some(connChnl.ChnlID),
			})
			return err
		})
		if getErr != nil {
			s.log.Error("step taking failed", compAttr)
			return execMod, execEff, connMod, getErr
		}
		connMod.ConnID = connSnap.ConnID
		connMod.ProviderCompID = connSnap.ProviderCompID
		connMod.ProviderCompRN = connSnap.ProviderCompRN
		connAttr := slog.Any("connID", connSnap.ConnID)
		subscription := connSnap.NextAct()
		if subscription != nil {
			connMod.OffsetNr = option.Some(subscription.GetActNr())
		}
		if subscription == nil {
			// регистрируем подписку доступопринимателя
			connMod.Acts = append(connMod.Acts, connact.SubRec{
				ConnID:  connSnap.ConnID,
				CompRef: execSnap.CompRef,
				ChnlID:  connChnl.ChnlID,
				ContExp: termexp.DetachRec{
					ContChnlPH: connChnl.ChnlPH,
				},
			})
			s.log.Debug("taking half done", compAttr, connAttr)
			return execMod, execEff, connMod, nil
		}
		releasing, ok := subscription.(connact.SubRec)
		if !ok {
			panic(connact.ErrActRecTypeUnexpected(subscription))
		}
		newChnlID := identity.New()
		execMod.Vars = append(execMod.Vars, compvar.LinearRec{
			CompRef: connChnl.CompRef,
			ConnID:  connChnl.ConnID,
			ChnlID:  newChnlID,
			ChnlPH:  connChnl.ChnlPH,
			ChnlBS:  connChnl.ChnlBS,
			ExpCK:   nextExpCK,
		})
		switch expRec := releasing.ContExp.(type) {
		case termexp.ReleaseRec:
			// вяжем продолжение доступовозвращателя
			execMod.Vars = append(execMod.Vars, compvar.LinearRec{
				CompRef: releasing.CompRef,
				ConnID:  releasing.ConnID,
				ChnlID:  newChnlID,
				ChnlPH:  expRec.ContChnlPH,
				ChnlBS:  connChnl.ChnlBS,
				ExpCK:   nextExpCK,
			})
		default:
			panic(termexp.ErrRecTypeUnexpected(releasing.ContExp))
		}
		s.log.Debug("step taking succeed", compAttr, connAttr)
		return execMod, execEff, connMod, nil
	default:
		panic(termexp.ErrSpecTypeUnexpected(exp))
	}
}

func (s *service) retrieveSnap(ref compsem.SemRef) (_ ExecSnap3, err error) {
	ctx := context.Background()
	var execSnap ExecSnap2
	getErr1 := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		execSnap, err = s.compExecRepo.GetSnapByRef(txCtx, conn, ref)
		return err
	})
	if getErr1 != nil {
		return ExecSnap3{}, getErr1
	}
	var structExps map[symbol.ADT]typeexp.ExpRec
	var linearExps map[symbol.ADT]typeexp.ExpRec
	getErr2 := s.transactor.ImplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		structExps, err = s.typeExpRepo.GetRecMap(txCtx, conn, ExtractExpCKs(execSnap.StructVars))
		if err != nil {
			return err
		}
		linearExps, err = s.typeExpRepo.GetRecMap(txCtx, conn, ExtractExpCKs(execSnap.LinearVars))
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
