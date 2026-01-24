package poolexec

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-runtime/lib/db"

	"orglang/go-runtime/adt/identity"
	"orglang/go-runtime/adt/procdec"
	"orglang/go-runtime/adt/procexec"
	"orglang/go-runtime/adt/revnum"
	"orglang/go-runtime/adt/typedef"
	"orglang/go-runtime/adt/typeexp"
	"orglang/go-runtime/adt/uniqsym"
)

// Port
type API interface {
	Run(ExecSpec) (ExecRef, error) // aka Create
	RetrieveSnap(identity.ADT) (ExecSnap, error)
	RetreiveRefs() ([]ExecRef, error)
	Spawn(procexec.ExecSpec) (procexec.ExecRef, error)
	Poll(PollSpec) (procexec.ExecRef, error)
}

type ExecSpec struct {
	PoolQN uniqsym.ADT
	SupID  identity.ADT
}

type ExecRef struct {
	ExecID identity.ADT
	ProcID identity.ADT // main
}

type ExecRec struct {
	ExecID identity.ADT
	ProcID identity.ADT // main
	SupID  identity.ADT
	ExecRN revnum.ADT
}

type ExecSnap struct {
	ExecID   identity.ADT
	Title    string
	SubExecs []ExecRef
}

type PollSpec struct {
	ExecID identity.ADT
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
	return &service{}
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

func (s *service) Run(spec ExecSpec) (ExecRef, error) {
	ctx := context.Background()
	s.log.Debug("creation started", slog.Any("spec", spec))
	execRec := ExecRec{
		ExecID: identity.New(),
		ProcID: identity.New(),
		SupID:  spec.SupID,
		ExecRN: revnum.New(),
	}
	liab := procexec.Liab{
		PoolID:  execRec.ExecID,
		PoolRN:  execRec.ExecRN,
		ExecRef: procexec.ExecRef{ID: execRec.ProcID},
	}
	err := s.operator.Explicit(ctx, func(ds db.Source) error {
		err := s.poolExecs.Insert(ds, execRec)
		if err != nil {
			s.log.Error("creation failed")
			return err
		}
		err = s.poolExecs.InsertLiab(ds, liab)
		if err != nil {
			s.log.Error("creation failed")
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("creation failed")
		return ExecRef{}, err
	}
	s.log.Debug("creation succeed", slog.Any("poolID", execRec.ExecID))
	return ConvertRecToRef(execRec), nil
}

func (s *service) Poll(spec PollSpec) (procexec.ExecRef, error) {
	return procexec.ExecRef{}, nil
}

func (s *service) Spawn(spec procexec.ExecSpec) (_ procexec.ExecRef, err error) {
	procAttr := slog.Any("procID", spec.ExecID)
	s.log.Debug("spawning started", procAttr)
	return procexec.ExecRef{}, nil
}

func (s *service) RetrieveSnap(execID identity.ADT) (snap ExecSnap, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		snap, err = s.poolExecs.SelectSubs(ds, execID)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("execID", execID))
		return ExecSnap{}, err
	}
	return snap, nil
}

func (s *service) RetreiveRefs() (refs []ExecRef, err error) {
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
