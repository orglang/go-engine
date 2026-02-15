package procdec

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/procbind"
	"orglang/go-engine/adt/synonym"
	"orglang/go-engine/adt/uniqref"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

type API interface {
	Incept(uniqsym.ADT) (DecRef, error)
	Create(DecSpec) (DecRef, error)
	RetrieveSnap(DecRef) (DecSnap, error)
	RetreiveRefs() ([]DecRef, error)
}

type DecRef = uniqref.ADT

type DecSpec struct {
	ProcQN uniqsym.ADT
	// endpoint where process acts as a provider
	ProviderBS procbind.BindSpec
	// endpoints where process acts as a client
	ClientBSes []procbind.BindSpec
}

type DecRec struct {
	DecRef     DecRef
	SynVK      valkey.ADT
	ProviderBS procbind.BindSpec
	ClientBSes []procbind.BindSpec
}

// aka ExpDec or ExpDecDef without expression
type DecSnap struct {
	DecRef     DecRef
	ProviderBS procbind.BindSpec
	ClientBSes []procbind.BindSpec
}

type service struct {
	procDecs Repo
	synonyms synonym.Repo
	operator db.Operator
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(procDecs Repo, synDecs synonym.Repo, operator db.Operator, log *slog.Logger) *service {
	return &service{procDecs, synDecs, operator, log}
}

func (s *service) Incept(procQN uniqsym.ADT) (_ DecRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("procQN", procQN)
	s.log.Debug("inception started", qnAttr)
	synVK, err := procQN.Key()
	if err != nil {
		return DecRef{}, err
	}
	newSyn := synonym.Rec{SynQN: procQN, SynVK: synVK}
	newDec := DecRec{DecRef: uniqref.New(), SynVK: synVK}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.synonyms.InsertRec(ds, newSyn)
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
		return DecRef{}, err
	}
	s.log.Debug("inception succeed", qnAttr, slog.Any("decRef", newDec.DecRef))
	return newDec.DecRef, nil
}

func (s *service) Create(spec DecSpec) (_ DecRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("procQN", spec.ProcQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newRec := DecRec{
		DecRef:     uniqref.New(),
		ProviderBS: spec.ProviderBS,
		ClientBSes: spec.ClientBSes,
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
		return DecRef{}, err
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("decRef", newRec.DecRef))
	return newRec.DecRef, nil
}

func (s *service) RetrieveSnap(ref DecRef) (snap DecSnap, err error) {
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

func (s *service) RetreiveRefs() (refs []DecRef, err error) {
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
		typeQNs = append(typeQNs, rec.ProviderBS.TypeQN)
		for _, y := range rec.ClientBSes {
			typeQNs = append(typeQNs, y.TypeQN)
		}
	}
	return typeQNs
}

func ErrRootMissingInEnv(rid identity.ADT) error {
	return fmt.Errorf("root missing in env: %v", rid)
}
