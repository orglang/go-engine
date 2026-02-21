package poolexec

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implsubst"
	"orglang/go-engine/adt/pooldec"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/procdec"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/typedef"
	"orglang/go-engine/adt/typeexp"
	"orglang/go-engine/adt/uniqsym"
)

type API interface {
	Run(ExecSpec) (implsem.SemRef, error) // aka Create
	RetrieveSnap(implsem.SemRef) (ExecSnap, error)
	RetreiveRefs() ([]implsem.SemRef, error)
	Take(poolstep.StepSpec) error
	Poll(PollSpec) (implsem.SemRef, error)
}

type ExecSpec struct {
	// ссылка на декларацию пула
	DescQN uniqsym.ADT
	// внутренняя и внешняя ссылки на вновь создаваемый пул
	ProviderSS implsubst.SubstSpec
	// подстановки ранее созданных пулов
	ClientSSes []implsubst.SubstSpec
}

type ExecRec struct {
	ImplRef    implsem.SemRef
	ProviderPH symbol.ADT
	ClientSRs  []implsubst.SubstRec
}

type ExecSnap struct {
	ExecRef  implsem.SemRef
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
	poolExecs Repo,
	implSems implsem.Repo,
	poolDecs pooldec.Repo,
	procDecs procdec.Repo,
	typeDefs typedef.Repo,
	typeExps typeexp.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{poolExecs, implSems, poolDecs, procDecs, typeDefs, typeExps, operator, log.With(name)}
}

func (s *service) Run(spec ExecSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	ssAttr := slog.Any("ss", spec.ProviderSS)
	s.log.Debug("starting creation...", ssAttr, slog.Any("spec", spec))
	execQNs := make([]uniqsym.ADT, 0, len(spec.ClientSSes)+1)
	for _, ss := range spec.ClientSSes {
		if ss.ImplQN == spec.ProviderSS.ImplQN {
			continue
		}
		execQNs = append(execQNs, ss.ImplQN)
	}
	var execRefs map[uniqsym.ADT]implsem.SemRef
	selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		execRefs, err = s.implSems.SelectRefsByQNs(ds, execQNs)
		return err
	})
	if selectErr != nil {
		s.log.Error("creation failed", ssAttr)
		return implsem.SemRef{}, selectErr
	}
	newRef := implsem.NewRef()
	newBind := implsem.SemBind{ImplQN: spec.ProviderSS.ImplQN, ImplID: newRef.ImplID}
	newImpl := implsem.SemRec{Ref: newRef, Bind: newBind, Kind: implsem.Pool}
	clientSRs := make([]implsubst.SubstRec, 0, len(spec.ClientSSes))
	for _, ss := range spec.ClientSSes {
		if ss.ImplQN == spec.ProviderSS.ImplQN {
			clientSRs = append(clientSRs, implsubst.SubstRec{ChnlPH: ss.ChnlPH, ImplID: newRef.ImplID})
			continue
		}
		clientSRs = append(clientSRs, implsubst.SubstRec{ChnlPH: ss.ChnlPH, ImplID: execRefs[ss.ImplQN].ImplID})
	}
	newExec := ExecRec{ImplRef: newRef, ProviderPH: spec.ProviderSS.ChnlPH, ClientSRs: clientSRs}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSems.InsertRec(ds, newImpl)
		if err != nil {
			return err
		}
		return s.poolExecs.InsertRec(ds, newExec)
	})
	if transactErr != nil {
		s.log.Error("creation failed", ssAttr)
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", ssAttr, slog.Any("ref", newRef))
	return newRef, nil
}

func (s *service) Poll(spec PollSpec) (implsem.SemRef, error) {
	return implsem.SemRef{}, nil
}

func (s *service) Take(spec poolstep.StepSpec) (err error) {
	refAttr := slog.Any("ref", spec.ImplRef)
	s.log.Debug("starting taking...", refAttr)
	return nil
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

func (s *service) RetreiveRefs() (refs []implsem.SemRef, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		refs, err = s.poolExecs.SelectRefs(ds)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}
