package poolexec

import (
	"context"
	"log/slog"
	"maps"
	"reflect"
	"slices"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/pooldec"
	"orglang/go-engine/adt/poolvar"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Run(ExecSpec) (implsem.SemRef, error) // aka Create
	RetrieveSnap(implsem.SemRef) (ExecSnap, error)
	Poll(PollSpec) (implsem.SemRef, error)
}

type ExecSpec struct {
	// ссылка на декларацию пула
	DescQN uniqsym.ADT
	// внутренняя и внешняя ссылки на вновь создаваемый пул
	ProviderVS implvar.VarSpec
	// внутренние и внешние ссылки на ранее созданные пулы
	ClientVSes []implvar.VarSpec
}

type ExecRec implvar.VarRec

type ExecSnap struct {
	ImplRef  implsem.SemRef
	Title    string
	SubExecs []implsem.SemRef
}

type PollSpec struct {
	ExecID identity.ADT
}

// ответственность за процесс
type Liab struct {
	// позитивное значение при вручении
	// негативное значение при лишении
	ExecRef implsem.SemRef
	ProcID  identity.ADT
}

type service struct {
	poolExecs Repo
	implSems  implsem.Repo
	poolDecs  pooldec.Repo
	poolVars  poolvar.Repo
	xactDefs  xactexp.Repo
	xactExps  xactexp.Repo
	operator  db.Operator
	log       *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	poolExecs Repo,
	implSems implsem.Repo,
	poolDecs pooldec.Repo,
	poolVars poolvar.Repo,
	xactDefs xactexp.Repo,
	xactExps xactexp.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{poolExecs, implSems, poolDecs, poolVars, xactDefs, xactExps, operator, log.With(name)}
}

func (s *service) Run(spec ExecSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	vsAttr := slog.Any("ss", spec.ProviderVS)
	s.log.Debug("starting creation...", vsAttr, slog.Any("spec", spec))
	// нужно выбрать провайдерские переменные клиентов!
	execQNs := make([]uniqsym.ADT, 0, len(spec.ClientVSes))
	for _, varSpec := range spec.ClientVSes {
		if varSpec.ImplQN == spec.ProviderVS.ImplQN {
			continue
		}
		execQNs = append(execQNs, varSpec.ImplQN)
	}
	var execRecs map[uniqsym.ADT]ExecRec
	selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		execRecs, err = s.poolExecs.SelectRecsByQNs(ds, execQNs)
		return err
	})
	if selectErr != nil {
		s.log.Error("creation failed", vsAttr)
		return implsem.SemRef{}, selectErr
	}
	newRef := implsem.NewRef()
	newBind := implsem.SemBind{ImplQN: spec.ProviderVS.ImplQN, ImplID: newRef.ImplID}
	newImpl := implsem.SemRec{Ref: newRef, Bind: newBind, Kind: implsem.Pool}
	newExec := ExecRec{
		ImplRef: newRef,
		ChnlID:  identity.New(),
		ChnlBS:  implvar.Provider,
		ChnlPH:  spec.ProviderVS.ChnlPH,
		// TODO: заполнить ExpVK
	}
	newVars := make([]implvar.VarRec, 0, len(spec.ClientVSes))
	for _, varSpec := range spec.ClientVSes {
		var implRef implsem.SemRef
		var chnlID identity.ADT
		if varSpec.ImplQN == spec.ProviderVS.ImplQN {
			implRef = newExec.ImplRef
			chnlID = newExec.ChnlID
		} else {
			implRef = execRecs[varSpec.ImplQN].ImplRef
			chnlID = execRecs[varSpec.ImplQN].ChnlID
		}
		newVars = append(newVars, implvar.VarRec{
			ImplRef: implRef,
			ChnlID:  chnlID,
			ChnlPH:  varSpec.ChnlPH,
			ChnlBS:  implvar.Client,
		})
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSems.InsertRec(ds, newImpl)
		if err != nil {
			return err
		}
		err = s.poolExecs.InsertRec(ds, newExec)
		if err != nil {
			return err
		}
		err = s.poolVars.InsertRecs(ds, newVars)
		if err != nil {
			return err
		}
		return s.poolExecs.TouchRecs(ds, ConvertRecsToRefs(slices.Collect(maps.Values(execRecs))))
	})
	if transactErr != nil {
		s.log.Error("creation failed", vsAttr)
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", vsAttr, slog.Any("ref", newRef))
	return newRef, nil
}

func (s *service) Poll(spec PollSpec) (implsem.SemRef, error) {
	return implsem.SemRef{}, nil
}

func (s *service) RetrieveSnap(ref implsem.SemRef) (snap ExecSnap, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		snap, err = s.poolExecs.SelectSubs(ds, ref)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("ref", ref))
		return ExecSnap{}, err
	}
	return snap, nil
}
