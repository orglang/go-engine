package poolexec

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/procdec"
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
	PoolQN uniqsym.ADT
}

type ExecRec struct {
	ExecRef implsem.SemRef
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
	procDecs procdec.Repo,
	typeDefs typedef.Repo,
	typeExps typeexp.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{poolExecs, procDecs, typeDefs, typeExps, operator, log.With(name)}
}

func (s *service) Run(spec ExecSpec) (implsem.SemRef, error) {
	ctx := context.Background()
	s.log.Debug("creation started", slog.Any("spec", spec))
	execRec := ExecRec{
		ExecRef: implsem.NewRef(),
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		return s.poolExecs.InsertRec(ds, execRec)
	})
	if transactErr != nil {
		s.log.Error("creation failed")
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", slog.Any("execRef", execRec.ExecRef))
	return execRec.ExecRef, nil
}

func (s *service) Poll(spec PollSpec) (implsem.SemRef, error) {
	return implsem.SemRef{}, nil
}

func (s *service) Take(spec poolstep.StepSpec) (err error) {
	qnAttr := slog.Any("procQN", spec.ProcQN)
	s.log.Debug("spawning started", qnAttr)
	return nil
}

func (s *service) RetrieveSnap(ref implsem.SemRef) (snap ExecSnap, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		snap, err = s.poolExecs.SelectSubs(ds, ref)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("execID", ref))
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
