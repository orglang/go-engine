package procdec

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/descvar"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/uniqsym"
)

type API interface {
	Incept(uniqsym.ADT) (descsem.SemRef, error)
	Create(DecSpec) (DecSnap, error)
	RetrieveSnap(descsem.SemRef) (DecSnap, error)
	RetreiveRefs() ([]descsem.SemRef, error)
}

type DecSpec struct {
	DescQN uniqsym.ADT
	// endpoint where process acts as a provider
	ProviderVS descvar.VarSpec
	// endpoints where process acts as a client
	ClientVSes []descvar.VarSpec
}

type DecRec struct {
	DescRef    descsem.SemRef
	ProviderVR descvar.VarRec
	ClientVRs  []descvar.VarRec
}

// aka ExpDec or ExpDecDef without expression
type DecSnap struct {
	DescRef    descsem.SemRef
	ProviderVR descvar.VarRec
	ClientVRs  []descvar.VarRec
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
	s.log.Debug("starting inception...", qnAttr)
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
	s.log.Debug("inception succeed", qnAttr, slog.Any("ref", newDec.DescRef))
	return newDec.DescRef, nil
}

func (s *service) Create(spec DecSpec) (_ DecSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", spec.DescQN)
	s.log.Debug("starting creation...", qnAttr, slog.Any("spec", spec))
	typeQNs := make([]uniqsym.ADT, 0, len(spec.ClientVSes)+1)
	for _, spec := range spec.ClientVSes {
		typeQNs = append(typeQNs, spec.DescQN)
	}
	var typeRefs map[uniqsym.ADT]descsem.SemRef
	selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		typeRefs, err = s.descSems.SelectRefsByQNs(ds, append(typeQNs, spec.ProviderVS.DescQN))
		return err
	})
	if selectErr != nil {
		return DecSnap{}, selectErr
	}
	providerVR := descvar.VarRec{
		ChnlPH: spec.ProviderVS.ChnlPH,
		DescID: typeRefs[spec.ProviderVS.DescQN].DescID,
	}
	clientVRs := make([]descvar.VarRec, 0, len(spec.ClientVSes))
	for _, vs := range spec.ClientVSes {
		clientVRs = append(clientVRs, descvar.VarRec{ChnlPH: vs.ChnlPH, DescID: typeRefs[vs.DescQN].DescID})
	}
	newRef := descsem.NewRef()
	newBind := descsem.SemBind{DescQN: spec.DescQN, DescID: newRef.DescID}
	newDesc := descsem.SemRec{Ref: newRef, Bind: newBind, Kind: descsem.Proc}
	newDec := DecRec{DescRef: newRef, ProviderVR: providerVR, ClientVRs: clientVRs}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descSems.InsertRec(ds, newDesc)
		if err != nil {
			return err
		}
		return s.procDecs.InsertRec(ds, newDec)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return DecSnap{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newRef))
	return DecSnap{DescRef: newRef, ProviderVR: providerVR, ClientVRs: clientVRs}, nil
}

func (s *service) RetrieveSnap(ref descsem.SemRef) (snap DecSnap, err error) {
	ctx := context.Background()
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		snap, err = s.procDecs.SelectSnap(ds, ref)
		return err
	})
	if err != nil {
		s.log.Error("retrieval failed", slog.Any("ref", ref))
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
	return []uniqsym.ADT{}
}

func ErrRootMissingInEnv(rid identity.ADT) error {
	return fmt.Errorf("root missing in env: %v", rid)
}
