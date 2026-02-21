package procdec

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/uniqsym"
)

type API interface {
	Incept(uniqsym.ADT) (descsem.SemRef, error)
	Create(DecSpec) (descsem.SemRef, error)
	RetrieveSnap(descsem.SemRef) (DecSnap, error)
	RetreiveRefs() ([]descsem.SemRef, error)
}

type DecSpec struct {
	ProcQN uniqsym.ADT
	// endpoint where process acts as a provider
	ProviderVS implvar.VarSpec
	// endpoints where process acts as a client
	ClientVSes []implvar.VarSpec
}

type DecRec struct {
	DescRef    descsem.SemRef
	ProviderVS implvar.VarSpec
	ClientVSes []implvar.VarSpec
}

// aka ExpDec or ExpDecDef without expression
type DecSnap struct {
	DescRef    descsem.SemRef
	ProviderVS implvar.VarSpec
	ClientVSes []implvar.VarSpec
}

type service struct {
	procDecs Repo
	descSems descsem.Repo
	operator db.Operator
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(procDecs Repo, descSems descsem.Repo, operator db.Operator, log *slog.Logger) *service {
	return &service{procDecs, descSems, operator, log}
}

func (s *service) Incept(procQN uniqsym.ADT) (_ descsem.SemRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", procQN)
	s.log.Debug("inception started", qnAttr)
	newRef := descsem.NewRef()
	newDec := DecRec{DescRef: newRef}
	newBind := descsem.SemBind{DescQN: procQN, DescID: newRef.DescID}
	newDesc := descsem.SemRec{Ref: newRef, Bind: newBind, Kind: descsem.Proc}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descSems.InsertRec(ds, newDesc)
		if err != nil {
			return err
		}
		err = s.procDecs.InsertRec(ds, newDec)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("inception failed", qnAttr)
		return descsem.SemRef{}, err
	}
	s.log.Debug("inception succeed", qnAttr, slog.Any("decRef", newDec.DescRef))
	return newDec.DescRef, nil
}

func (s *service) Create(spec DecSpec) (_ descsem.SemRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("procQN", spec.ProcQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newRec := DecRec{
		DescRef:    descsem.NewRef(),
		ProviderVS: spec.ProviderVS,
		ClientVSes: spec.ClientVSes,
	}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.procDecs.InsertRec(ds, newRec)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("creation failed", qnAttr)
		return descsem.SemRef{}, err
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("decRef", newRec.DescRef))
	return newRec.DescRef, nil
}

func (s *service) RetrieveSnap(ref descsem.SemRef) (snap DecSnap, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		snap, err = s.procDecs.SelectSnap(ds, ref)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("decRef", ref))
		return DecSnap{}, err
	}
	return snap, nil
}

func (s *service) RetreiveRefs() (refs []descsem.SemRef, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		refs, err = s.procDecs.SelectRefs(ds)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed")
		return nil, err
	}
	return refs, nil
}

func CollectEnv(recs iter.Seq[DecRec]) []uniqsym.ADT {
	typeQNs := []uniqsym.ADT{}
	for rec := range recs {
		typeQNs = append(typeQNs, rec.ProviderVS.TypeQN)
		for _, y := range rec.ClientVSes {
			typeQNs = append(typeQNs, y.TypeQN)
		}
	}
	return typeQNs
}

func ErrRootMissingInEnv(rid identity.ADT) error {
	return fmt.Errorf("root missing in env: %v", rid)
}
