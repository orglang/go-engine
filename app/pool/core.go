package pool

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/exp/maps"

	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
	"smecalculus/rolevod/lib/rev"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/proc"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/role"
	"smecalculus/rolevod/app/sig"
)

type ID = id.ADT
type Rev = rev.ADT
type Title = string

type Spec struct {
	Title  string
	SupID  id.ADT
	DepIDs []sig.ID
}

type Ref struct {
	PoolID id.ADT
	Title  string
}

type Mutex struct {
	PoolID id.ADT
	Kind   rev.Knd
	Rev    rev.ADT
}

type Root struct {
	PoolID id.ADT
	Title  string
	SupID  id.ADT
	Revs   []rev.ADT
}

type Transfer struct {
	ProcID   id.ADT
	PoolID   id.ADT
	ExPoolID id.ADT
	Rev      rev.ADT
}

const (
	rootLock rev.Knd = iota
	procLock
)

type Mod struct {
	Mutexes []Mutex
	Assets  []Transfer
}

type RootMod struct {
	PoolID id.ADT
	Rev    rev.ADT
	K      rev.Knd
}

type SubSnap struct {
	PoolID id.ADT
	Title  string
	Subs   []Ref
}

type AssetSnap struct {
	PoolID id.ADT
	Title  string
	EPs    []proc.Chnl
}

type AssetMod struct {
	OutPoolID id.ADT
	InPoolID  id.ADT
	Rev       rev.ADT
	EPs       []proc.Chnl
}

type LiabSnap struct {
	PoolID id.ADT
	Title  string
	EP     proc.Chnl
}

type LiabMod struct {
	OutPoolID id.ADT
	InPoolID  id.ADT
	Rev       rev.ADT
	EP        proc.Chnl
}

type TranSpec struct {
	PoolID id.ADT
	ProcID id.ADT
	Term   step.Term
}

type Environment struct {
	sigs   map[sig.ID]sig.Root
	roles  map[role.QN]role.Root
	states map[state.ID]state.Root
}

func (e Environment) Contains(id sig.ID) bool {
	_, ok := e.sigs[id]
	return ok
}

func (e Environment) LookupPE(id sig.ID) state.EP {
	decl := e.sigs[id]
	role := e.roles[decl.PE.Link]
	return state.EP{Z: decl.PE.Link, C: e.states[role.StateID]}
}

func (e Environment) LookupCEs(id sig.ID) []state.EP {
	decl := e.sigs[id]
	ces := []state.EP{}
	for _, ce := range decl.CEs {
		role := e.roles[decl.PE.Link]
		ces = append(ces, state.EP{Z: ce.Link, C: e.states[role.StateID]})
	}
	return ces
}

// Port
type API interface {
	Create(Spec) (Root, error)
	Retrieve(id.ADT) (SubSnap, error)
	RetreiveRefs() ([]Ref, error)
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

type service struct {
	pools    Repo
	sigs     sig.Repo
	roles    role.Repo
	states   state.Repo
	operator data.Operator
	log      *slog.Logger
}

func newService(
	pools Repo,
	sigs sig.Repo,
	roles role.Repo,
	states state.Repo,
	operator data.Operator,
	l *slog.Logger,
) *service {
	name := slog.String("name", "poolService")
	return &service{pools, sigs, roles, states, operator, l.With(name)}
}

func (s *service) Create(spec Spec) (_ Root, err error) {
	ctx := context.Background()
	s.log.Debug("creation started", slog.Any("spec", spec))
	root := Root{
		PoolID: id.New(),
		Revs:   []rev.ADT{rev.Initial()},
		Title:  spec.Title,
		SupID:  spec.SupID,
	}
	s.operator.Explicit(ctx, func(ds data.Source) error {
		err = s.pools.Insert(ds, root)
		return err
	})
	if err != nil {
		s.log.Error("creation failed")
		return Root{}, err
	}
	s.log.Debug("creation succeeded", slog.Any("id", root.PoolID))
	return root, nil
}

func (s *service) Take(spec TranSpec) (err error) {
	ctx := context.Background()
	// initial values
	poolID := spec.PoolID
	procID := spec.ProcID
	term := spec.Term
	for term != nil {
		var procSnap proc.Snap
		s.operator.Implicit(ctx, func(ds data.Source) {
			procSnap, err = s.pools.SelectProc(ds, procID)
		})
		idAttr := slog.Any("procID", procID)
		if err != nil {
			s.log.Error("taking failed", idAttr)
			return err
		}
		if len(procSnap.Chnls) == 0 {
			s.log.Error("taking failed", idAttr)
			return err
		}
		sigIDs := step.CollectEnv(term)
		var sigs map[sig.ID]sig.Root
		s.operator.Implicit(ctx, func(ds data.Source) {
			sigs, err = s.sigs.SelectEnv(ds, sigIDs)
		})
		if err != nil {
			s.log.Error("taking failed", idAttr, slog.Any("sigs", sigIDs))
			return err
		}
		roleQNs := sig.CollectEnv(maps.Values(sigs))
		var roles map[role.QN]role.Root
		s.operator.Implicit(ctx, func(ds data.Source) {
			roles, err = s.roles.SelectEnv(ds, roleQNs)
		})
		if err != nil {
			s.log.Error("taking failed", idAttr, slog.Any("roles", roleQNs))
			return err
		}
		envIDs := role.CollectEnv(maps.Values(roles))
		ctxIDs := CollectCtx(maps.Values(procSnap.Chnls))
		var states map[state.ID]state.Root
		s.operator.Implicit(ctx, func(ds data.Source) {
			states, err = s.states.SelectEnv(ds, append(envIDs, ctxIDs...))
		})
		if err != nil {
			s.log.Error("taking failed", idAttr, slog.Any("env", envIDs), slog.Any("ctx", ctxIDs))
			return err
		}
		environ := Environment{sigs, roles, states}
		context := convertToCtx(maps.Values(procSnap.Chnls), states)
		// type checking
		err = s.checkState(poolID, environ, context, procSnap, term)
		if err != nil {
			s.log.Error("taking failed", idAttr)
			return err
		}
		// step taking
		nextSpec, procMod, err := s.takeWith(context, procSnap, term)
		if err != nil {
			s.log.Error("taking failed", idAttr)
			return err
		}
		s.operator.Explicit(ctx, func(ds data.Source) error {
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
		term = nextSpec.Term
	}
	return nil
}

func (s *service) takeWith(
	context state.Context,
	snap proc.Snap,
	t step.Term,
) (
	spec TranSpec,
	mod proc.Mod,
	_ error,
) {
	switch termSpec := t.(type) {
	case step.CloseSpec:
		viaChnl, ok := snap.Chnls[termSpec.X]
		if !ok {
			err := chnl.ErrMissingInCfg(termSpec.X)
			s.log.Error("taking failed")
			return spec, mod, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		lock := proc.Lock{
			PoolID: viaChnl.PS.PoolID,
			Rev:    viaChnl.PS.Rev,
		}
		mod.Locks = append(mod.Locks, lock)
		viaStep, ok := snap.Steps[viaChnl.ChnlID]
		if !ok {
			err := step.ErrMissingInCfg(viaChnl.ChnlID)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		sndrViaBnd := proc.Bnd{
			ChnlPH:  termSpec.X,
			ChnlID:  id.Nil,
			StateID: id.Nil,
			ProcID:  snap.ProcID,
			Rev:     snap.Rev + 1,
		}
		mod.Bnds = append(mod.Bnds, sndrViaBnd)
		if viaStep == nil {
			msgStep := step.MsgRoot2{
				PoolID: snap.PoolID,
				ProcID: snap.ProcID,
				ChnlID: viaChnl.ChnlID,
				Rev:    snap.Rev,
				Val:    step.CloseImpl{},
			}
			mod.Steps = append(mod.Steps, msgStep)
			s.log.Debug("taking half done", viaAttr)
			return spec, mod, nil
		}
		svcStep, ok := viaStep.(step.SvcRoot2)
		if !ok {
			err := step.ErrRootTypeUnexpected(viaStep)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		waitImpl, ok := svcStep.Cont.(step.WaitImpl)
		if !ok {
			err := fmt.Errorf("cont type mismatch: want %T, got %T", waitImpl, svcStep.Cont)
			s.log.Error("taking failed", viaAttr, slog.Any("cont", svcStep.Cont))
			return spec, mod, err
		}
		spec = TranSpec{
			PoolID: svcStep.PoolID,
			ProcID: svcStep.ProcID,
			Term:   waitImpl.Cont,
		}
		s.log.Debug("taking succeeded", viaAttr)
		return spec, mod, nil
	case step.WaitSpec:
		viaChnl, ok := snap.Chnls[termSpec.X]
		if !ok {
			err := chnl.ErrMissingInCfg(termSpec.X)
			s.log.Error("taking failed")
			return spec, mod, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		lock := proc.Lock{
			PoolID: viaChnl.CS.PoolID,
			Rev:    viaChnl.CS.Rev,
		}
		mod.Locks = append(mod.Locks, lock)
		viaStep, ok := snap.Steps[viaChnl.ChnlID]
		if !ok {
			err := step.ErrMissingInCfg(viaChnl.ChnlID)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		rcvrViaBnd := proc.Bnd{
			ChnlPH:  termSpec.X,
			ChnlID:  id.Nil,
			StateID: id.Nil,
			ProcID:  snap.ProcID,
			Rev:     snap.Rev + 1,
		}
		mod.Bnds = append(mod.Bnds, rcvrViaBnd)
		if viaStep == nil {
			svcStep := step.SvcRoot2{
				PoolID: snap.PoolID,
				ProcID: snap.ProcID,
				ChnlID: viaChnl.ChnlID,
				Cont: step.WaitImpl{
					Cont: termSpec.Cont,
				},
			}
			mod.Steps = append(mod.Steps, svcStep)
			s.log.Debug("taking half done", viaAttr)
			return spec, mod, nil
		}
		msgStep, ok := viaStep.(step.MsgRoot2)
		if !ok {
			err := step.ErrRootTypeUnexpected(viaStep)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		switch stepImpl := msgStep.Val.(type) {
		case step.CloseImpl:
			// all done
		case step.FwdImpl:
			viaState, ok := context.Linear[termSpec.X]
			if !ok {
				err := state.ErrMissingInCtx(termSpec.X)
				s.log.Error("taking failed", viaAttr)
				return spec, mod, err
			}
			fwdBnd := proc.Bnd{
				ChnlPH:  termSpec.X,
				ChnlID:  stepImpl.ChnlID,
				StateID: viaState.Ident(),
				ProcID:  viaChnl.CS.ProcID,
				Rev:     viaChnl.CS.Rev + 1,
			}
			mod.Bnds = append(mod.Bnds, fwdBnd)
			spec = TranSpec{
				PoolID: viaChnl.CS.PoolID,
				ProcID: viaChnl.CS.ProcID,
				Term:   termSpec,
			}
			s.log.Debug("taking succeeded", viaAttr)
			return spec, mod, nil
		default:
			panic(step.ErrValTypeUnexpected2(msgStep.Val))
		}
		spec = TranSpec{
			PoolID: viaChnl.CS.PoolID,
			ProcID: viaChnl.CS.ProcID,
			Term:   termSpec.Cont,
		}
		s.log.Debug("taking succeeded", viaAttr)
		return spec, mod, nil
	case step.SendSpec:
		viaChnl, ok := snap.Chnls[termSpec.X]
		if !ok {
			err := chnl.ErrMissingInCfg(termSpec.X)
			s.log.Error("taking failed")
			return spec, mod, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		lock := proc.Lock{
			PoolID: viaChnl.PS.PoolID,
			Rev:    viaChnl.PS.Rev,
		}
		mod.Locks = append(mod.Locks, lock)
		viaStep, ok := snap.Steps[viaChnl.ChnlID]
		if !ok {
			err := step.ErrMissingInCfg(viaChnl.ChnlID)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		valChnl, ok := snap.Chnls[termSpec.Y]
		if !ok {
			err := chnl.ErrMissingInCfg(termSpec.Y)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		sndrValBnd := proc.Bnd{
			ChnlPH:  termSpec.Y,
			ChnlID:  id.Nil,
			StateID: id.Nil,
			ProcID:  snap.ProcID,
			Rev:     snap.Rev + 1,
		}
		mod.Bnds = append(mod.Bnds, sndrValBnd)
		if viaStep == nil {
			msgStep := step.MsgRoot2{
				PoolID: snap.PoolID,
				ProcID: snap.ProcID,
				ChnlID: viaChnl.ChnlID,
				Rev:    snap.Rev,
				Val: step.SendImpl{
					X: termSpec.X,
					A: id.New(),
					B: valChnl.ChnlID,
				},
			}
			mod.Steps = append(mod.Steps, msgStep)
			s.log.Debug("taking half done", viaAttr)
			return spec, mod, nil
		}
		svcStep, ok := viaStep.(step.SvcRoot2)
		if !ok {
			err := step.ErrRootTypeUnexpected(viaStep)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		viaState, ok := context.Linear[termSpec.X]
		if !ok {
			err := state.ErrMissingInCtx(termSpec.X)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		viaStateID := viaState.(state.Prod).Next()
		switch termImpl := svcStep.Cont.(type) {
		case step.RecvImpl:
			rcvrViaBnd := proc.Bnd{
				ChnlPH:  termImpl.X,
				ChnlID:  termImpl.A,
				StateID: viaStateID,
				ProcID:  svcStep.ProcID,
				Rev:     svcStep.Rev + 1,
			}
			mod.Bnds = append(mod.Bnds, rcvrViaBnd)
			sndrViaBnd := proc.Bnd{
				ChnlPH:  termSpec.X,
				ChnlID:  termImpl.A,
				StateID: viaStateID,
				ProcID:  snap.ProcID,
				Rev:     snap.Rev + 1,
			}
			mod.Bnds = append(mod.Bnds, sndrViaBnd)
			rcvrValBnd := proc.Bnd{
				ChnlPH:  termImpl.Y,
				ChnlID:  valChnl.ChnlID,
				StateID: valChnl.StateID,
				ProcID:  svcStep.ProcID,
				Rev:     svcStep.Rev + 1,
			}
			mod.Bnds = append(mod.Bnds, rcvrValBnd)
			spec = TranSpec{
				PoolID: svcStep.PoolID,
				ProcID: svcStep.ProcID,
				Term:   termImpl.Cont,
			}
			s.log.Debug("taking succeeded", viaAttr)
			return spec, mod, nil
		default:
			panic(step.ErrContTypeUnexpected2(svcStep.Cont))
		}
	case step.RecvSpec:
		viaChnl, ok := snap.Chnls[termSpec.X]
		if !ok {
			err := chnl.ErrMissingInCfg(termSpec.X)
			s.log.Error("taking failed")
			return spec, mod, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		lock := proc.Lock{
			PoolID: viaChnl.CS.PoolID,
			Rev:    viaChnl.CS.Rev,
		}
		mod.Locks = append(mod.Locks, lock)
		viaStep, ok := snap.Steps[viaChnl.ChnlID]
		if !ok {
			err := step.ErrMissingInCfg(viaChnl.ChnlID)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		if viaStep == nil {
			svcStep := step.SvcRoot2{
				PoolID: snap.PoolID,
				ProcID: snap.ProcID,
				ChnlID: viaChnl.ChnlID,
				Rev:    snap.Rev,
				Cont: step.RecvImpl{
					X:    termSpec.X,
					A:    id.New(),
					Cont: termSpec.Cont,
				},
			}
			mod.Steps = append(mod.Steps, svcStep)
			s.log.Debug("taking half done", viaAttr)
			return spec, mod, nil
		}
		msgStep, ok := viaStep.(step.MsgRoot2)
		if !ok {
			err := step.ErrRootTypeUnexpected(viaStep)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		viaState, ok := context.Linear[termSpec.X]
		if !ok {
			err := state.ErrMissingInCtx(termSpec.X)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		viaStateID := viaState.(state.Prod).Next()
		valState, ok := context.Linear[termSpec.Y]
		if !ok {
			err := state.ErrMissingInCtx(termSpec.Y)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		switch termImpl := msgStep.Val.(type) {
		case step.SendImpl:
			sndrViaBnd := proc.Bnd{
				ChnlPH:  termImpl.X,
				ChnlID:  termImpl.A,
				StateID: viaStateID,
				ProcID:  msgStep.ProcID,
				Rev:     msgStep.Rev + 1,
			}
			mod.Bnds = append(mod.Bnds, sndrViaBnd)
			rcvrViaBnd := proc.Bnd{
				ChnlPH:  termSpec.X,
				ChnlID:  termImpl.A,
				StateID: viaStateID,
				ProcID:  snap.ProcID,
				Rev:     snap.Rev + 1,
			}
			mod.Bnds = append(mod.Bnds, rcvrViaBnd)
			rcvrValBnd := proc.Bnd{
				ChnlPH:  termSpec.Y,
				ChnlID:  termImpl.B,
				StateID: valState.Ident(),
				ProcID:  snap.ProcID,
				Rev:     snap.Rev + 1,
			}
			mod.Bnds = append(mod.Bnds, rcvrValBnd)
			spec = TranSpec{
				PoolID: snap.PoolID,
				ProcID: snap.ProcID,
				Term:   termSpec.Cont,
			}
			s.log.Debug("taking succeeded", viaAttr)
			return spec, mod, nil
		default:
			panic(step.ErrValTypeUnexpected2(msgStep.Val))
		}
	case step.LabSpec:
		viaChnl, ok := snap.Chnls[termSpec.X]
		if !ok {
			err := chnl.ErrMissingInCfg(termSpec.X)
			s.log.Error("taking failed")
			return spec, mod, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		lock := proc.Lock{
			PoolID: viaChnl.PS.PoolID,
			Rev:    viaChnl.PS.Rev,
		}
		mod.Locks = append(mod.Locks, lock)
		viaStep, ok := snap.Steps[viaChnl.ChnlID]
		if !ok {
			err := step.ErrMissingInCfg(viaChnl.ChnlID)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		if viaStep == nil {
			msgStep := step.MsgRoot2{
				PoolID: snap.PoolID,
				ProcID: snap.ProcID,
				ChnlID: viaChnl.ChnlID,
				Val: step.LabImpl{
					X: termSpec.X,
					A: id.New(),
					L: termSpec.L,
				},
			}
			mod.Steps = append(mod.Steps, msgStep)
			s.log.Debug("taking half done", viaAttr)
			return spec, mod, nil
		}
		svcStep, ok := viaStep.(step.SvcRoot2)
		if !ok {
			err := step.ErrRootTypeUnexpected(viaStep)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		termImpl, ok := svcStep.Cont.(step.CaseImpl)
		if !ok {
			err := fmt.Errorf("cont type mismatch: want %T, got %T", termImpl, svcStep.Cont)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		viaState, ok := context.Linear[termSpec.X]
		if !ok {
			err := state.ErrMissingInCtx(termSpec.X)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		viaStateID := viaState.(state.Sum).Next(termSpec.L)
		rcvrViaBnd := proc.Bnd{
			ChnlPH:  termImpl.X,
			ChnlID:  termImpl.A,
			StateID: viaStateID,
			ProcID:  svcStep.ProcID,
			Rev:     svcStep.Rev + 1,
		}
		mod.Bnds = append(mod.Bnds, rcvrViaBnd)
		sndrViaBnd := proc.Bnd{
			ChnlPH:  termSpec.X,
			ChnlID:  termImpl.A,
			StateID: viaStateID,
			ProcID:  snap.ProcID,
			Rev:     snap.Rev + 1,
		}
		mod.Bnds = append(mod.Bnds, sndrViaBnd)
		spec = TranSpec{
			PoolID: svcStep.PoolID,
			ProcID: svcStep.ProcID,
			Term:   termImpl.Conts[termSpec.L],
		}
		s.log.Debug("taking succeeded", viaAttr)
		return spec, mod, nil
	case step.CaseSpec:
		viaChnl, ok := snap.Chnls[termSpec.X]
		if !ok {
			err := chnl.ErrMissingInCfg(termSpec.X)
			s.log.Error("taking failed")
			return spec, mod, err
		}
		viaAttr := slog.Any("chnlID", viaChnl.ChnlID)
		lock := proc.Lock{
			PoolID: viaChnl.CS.PoolID,
			Rev:    viaChnl.CS.Rev,
		}
		mod.Locks = append(mod.Locks, lock)
		viaStep, ok := snap.Steps[viaChnl.ChnlID]
		if !ok {
			err := step.ErrMissingInCfg(viaChnl.ChnlID)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		if viaStep == nil {
			svcStep := step.SvcRoot2{
				PoolID: snap.PoolID,
				ProcID: snap.ProcID,
				ChnlID: viaChnl.ChnlID,
				Cont: step.CaseImpl{
					X:     termSpec.X,
					A:     id.New(),
					Conts: termSpec.Conts,
				},
			}
			mod.Steps = append(mod.Steps, svcStep)
			s.log.Debug("taking half done", viaAttr)
			return spec, mod, nil
		}
		msgStep, ok := viaStep.(step.MsgRoot2)
		if !ok {
			err := step.ErrRootTypeUnexpected(viaStep)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		termImpl := msgStep.Val.(step.LabImpl)
		if !ok {
			err := fmt.Errorf("cont type mismatch: want %T, got %T", termImpl, msgStep.Val)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		viaState, ok := context.Linear[termSpec.X]
		if !ok {
			err := state.ErrMissingInCtx(termSpec.X)
			s.log.Error("taking failed", viaAttr)
			return spec, mod, err
		}
		viaStateID := viaState.(state.Sum).Next(termImpl.L)
		sndrViaBnd := proc.Bnd{
			ChnlPH:  termImpl.X,
			ChnlID:  termImpl.A,
			StateID: viaStateID,
			ProcID:  msgStep.ProcID,
			Rev:     msgStep.Rev + 1,
		}
		mod.Bnds = append(mod.Bnds, sndrViaBnd)
		rcvrViaBnd := proc.Bnd{
			ChnlPH:  termSpec.X,
			ChnlID:  termImpl.A,
			StateID: viaStateID,
			ProcID:  snap.ProcID,
			Rev:     snap.Rev + 1,
		}
		mod.Bnds = append(mod.Bnds, rcvrViaBnd)
		spec = TranSpec{
			PoolID: snap.PoolID,
			ProcID: snap.ProcID,
			Term:   termSpec.Conts[termImpl.L],
		}
		s.log.Debug("taking succeeded", viaAttr)
		return spec, mod, nil
	default:
		panic(step.ErrTermTypeUnexpected(t))
	}
}

func (s *service) Retrieve(poolID id.ADT) (snap SubSnap, err error) {
	ctx := context.Background()
	s.operator.Implicit(ctx, func(ds data.Source) {
		snap, err = s.pools.SelectSubs(ds, poolID)
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("id", poolID))
		return SubSnap{}, err
	}
	return snap, nil
}

func (s *service) RetreiveRefs() (refs []Ref, err error) {
	ctx := context.Background()
	s.operator.Implicit(ctx, func(ds data.Source) {
		refs, err = s.pools.SelectRefs(ds)
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}

func CollectCtx(eps []proc.Chnl) []state.ID {
	return nil
}

func convertToCtx(eps []proc.Chnl, states map[state.ID]state.Root) state.Context {
	linear := make(map[ph.ADT]state.Root, len(eps))
	for _, ep := range eps {
		linear[ep.ChnlPH] = states[ep.StateID]
	}
	return state.Context{Linear: linear}
}

func convertToCfg(eps []proc.Chnl) map[ph.ADT]proc.Chnl {
	cfg := make(map[ph.ADT]proc.Chnl, len(eps))
	for _, ep := range eps {
		cfg[ep.ChnlPH] = ep
	}
	return cfg
}

func (s *service) checkState(
	poolID id.ADT,
	env Environment,
	ctx state.Context,
	acl proc.Snap,
	t step.Term,
) error {
	ch, ok := acl.Chnls[t.Via()]
	if !ok {
		panic("no via in acl")
	}
	if ch.SrvrID == ch.ClntID {
		panic("can not be equal")
	}
	switch poolID {
	case ch.SrvrID:
		return s.checkProvider(poolID, env, ctx, acl, t)
	case ch.ClntID:
		return s.checkClient(poolID, env, ctx, acl, t)
	default:
		s.log.Error("state checking failed", slog.Any("id", poolID))
		return fmt.Errorf("unknown pool id")
	}
}

func (s *service) checkProvider(
	poolID id.ADT,
	env Environment,
	ctx state.Context,
	acl proc.Snap,
	t step.Term,
) error {
	return nil
}

func (s *service) checkClient(
	poolID id.ADT,
	env Environment,
	ctx state.Context,
	acl proc.Snap,
	t step.Term,
) error {
	return nil
}

// Port
type Repo interface {
	Insert(data.Source, Root) error
	SelectRefs(data.Source) ([]Ref, error)
	SelectSubs(data.Source, id.ADT) (SubSnap, error)
	SelectAssets(data.Source, id.ADT) (AssetSnap, error)
	SelectProc(data.Source, id.ADT) (proc.Snap, error)
	UpdateProc(data.Source, proc.Mod) error
	UpdateAssets(data.Source, AssetMod) error
	Transfer(source data.Source, giver id.ADT, taker id.ADT, pids []chnl.ID) error
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	ConvertRootToRef func(Root) Ref
)

func errOptimisticUpdate(got rev.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
