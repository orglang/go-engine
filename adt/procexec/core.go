package procexec

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
	"orglang/orglang/adt/revnum"

	"orglang/orglang/adt/procdecl"
	"orglang/orglang/adt/procdef"
	"orglang/orglang/adt/typedef"
)

type API interface {
	Run(ProcSpec) error
	Retrieve(identity.ADT) (ProcSnap, error)
}

type SemRec interface {
	step() identity.ADT
}

func ChnlID(r SemRec) identity.ADT { return r.step() }

type MsgRec struct {
	PoolID identity.ADT
	ProcID identity.ADT
	ChnlID identity.ADT
	Val    procdef.TermRec
	PoolRN revnum.ADT
	ProcRN revnum.ADT
}

func (r MsgRec) step() identity.ADT { return r.ChnlID }

type SvcRec struct {
	PoolID identity.ADT
	ProcID identity.ADT
	ChnlID identity.ADT
	Cont   procdef.TermRec
	PoolRN revnum.ADT
}

func (r SvcRec) step() identity.ADT { return r.ChnlID }

type ProcSpec struct {
	PoolID identity.ADT
	ExecID identity.ADT
	ProcTS procdef.TermSpec
}

type ProcRef struct {
	ExecID identity.ADT
}

type ProcSnap struct {
	ExecID identity.ADT
}

type MainCfg struct {
	ProcID identity.ADT
	Bnds   map[qualsym.ADT]EP2
	Acts   map[identity.ADT]SemRec
	PoolID identity.ADT
	ProcRN revnum.ADT
}

// aka Configuration
type Cfg struct {
	ProcID identity.ADT
	Chnls  map[qualsym.ADT]EP
	Steps  map[identity.ADT]SemRec
	PoolID identity.ADT
	PoolRN revnum.ADT
	ProcRN revnum.ADT
}

type Env struct {
	ProcSigs  map[identity.ADT]procdecl.ProcRec
	Types     map[qualsym.ADT]typedef.TypeRec
	TypeTerms map[identity.ADT]typedef.TermRec
	Locks     map[qualsym.ADT]Lock
}

type EP struct {
	ChnlPH qualsym.ADT
	ChnlID identity.ADT
	TermID identity.ADT
	// provider
	PoolID identity.ADT
}

type EP2 struct {
	CordPH qualsym.ADT
	CordID identity.ADT
	// provider
	PoolID identity.ADT
}

type Lock struct {
	PoolID identity.ADT
	PoolRN revnum.ADT
}

func ChnlPH(rec EP) qualsym.ADT { return rec.ChnlPH }

// ответственность за процесс
type Liab struct {
	PoolID identity.ADT
	ProcID identity.ADT
	// позитивное значение при вручении
	// негативное значение при лишении
	PoolRN revnum.ADT
}

type Mod struct {
	Locks []Lock
	Bnds  []Bnd
	Steps []SemRec
	Liabs []Liab
}

type MainMod struct {
	Bnds []Bnd
	Acts []SemRec
}

type Bnd struct {
	ProcID identity.ADT
	ChnlPH qualsym.ADT
	ChnlID identity.ADT
	TermID identity.ADT
	PoolRN revnum.ADT
	ProcRN revnum.ADT
}

type service struct {
	procs    repo
	operator sd.Operator
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

func newService(
	procs repo,
	operator sd.Operator,
	l *slog.Logger,
) *service {
	return &service{procs, operator, l}
}

func (s *service) Run(spec ProcSpec) (err error) {
	idAttr := slog.Any("procID", spec.ExecID)
	s.log.Debug("creation started", idAttr)
	ctx := context.Background()
	var mainCfg MainCfg
	err = s.operator.Implicit(ctx, func(ds sd.Source) error {
		mainCfg, err = s.procs.SelectMain(ds, spec.ExecID)
		return err
	})
	if err != nil {
		s.log.Error("creation failed", idAttr)
		return err
	}
	var mainEnv Env
	err = s.checkType(spec.PoolID, mainEnv, mainCfg, spec.ProcTS)
	if err != nil {
		s.log.Error("creation failed", idAttr)
		return err
	}
	mainMod, err := s.createWith(mainEnv, mainCfg, spec.ProcTS)
	if err != nil {
		s.log.Error("creation failed", idAttr)
		return err
	}
	err = s.operator.Explicit(ctx, func(ds sd.Source) error {
		err = s.procs.UpdateMain(ds, mainMod)
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

func (s *service) Retrieve(procID identity.ADT) (_ ProcSnap, err error) {
	return ProcSnap{}, nil
}

func (s *service) checkType(
	poolID identity.ADT,
	mainEnv Env,
	mainCfg MainCfg,
	termSpec procdef.TermSpec,
) error {
	imp, ok := mainCfg.Bnds[termSpec.Via()]
	if !ok {
		panic("no via in main cfg")
	}
	if poolID == imp.PoolID {
		return s.checkProvider(poolID, mainEnv, mainCfg, termSpec)
	} else {
		return s.checkClient(poolID, mainEnv, mainCfg, termSpec)
	}
}

func (s *service) checkProvider(
	poolID identity.ADT,
	mainEnv Env,
	mainCfg MainCfg,
	ts procdef.TermSpec,
) error {
	return nil
}

func (s *service) checkClient(
	poolID identity.ADT,
	mainEnv Env,
	mainCfg MainCfg,
	ts procdef.TermSpec,
) error {
	return nil
}

func (s *service) createWith(
	mainEnv Env,
	procCfg MainCfg,
	ts procdef.TermSpec,
) (
	procMod MainMod,
	_ error,
) {
	switch termSpec := ts.(type) {
	case procdef.CallSpecOld:
		viaCord, ok := procCfg.Bnds[termSpec.X]
		if !ok {
			err := procdef.ErrMissingInCfg(termSpec.X)
			s.log.Error("coordination failed")
			return MainMod{}, err
		}
		viaAttr := slog.Any("cordID", viaCord.CordID)
		for _, chnlPH := range termSpec.Ys {
			sndrValBnd := Bnd{
				ProcID: procCfg.ProcID,
				ChnlPH: chnlPH,
				ProcRN: -procCfg.ProcRN.Next(),
			}
			procMod.Bnds = append(procMod.Bnds, sndrValBnd)
		}
		rcvrAct := procCfg.Acts[viaCord.CordID]
		if rcvrAct == nil {
			sndrAct := MsgRec{}
			procMod.Acts = append(procMod.Acts, sndrAct)
			s.log.Debug("coordination half done", viaAttr)
			return procMod, nil
		}
		s.log.Debug("coordination succeed")
		return procMod, nil
	case procdef.SpawnSpecOld:
		s.log.Debug("coordination succeed")
		return procMod, nil
	default:
		panic(procdef.ErrTermTypeUnexpected(ts))
	}
}

func ErrMissingChnl(want qualsym.ADT) error {
	return fmt.Errorf("channel missing in cfg: %v", want)
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

func ErrRootTypeUnexpected(got SemRec) error {
	return fmt.Errorf("sem rec unexpected: %T", got)
}

func ErrRootTypeMismatch(got, want SemRec) error {
	return fmt.Errorf("sem rec mismatch: want %T, got %T", want, got)
}
