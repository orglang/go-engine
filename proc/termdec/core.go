package termdec

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/semtype"
	"orglang/go-engine/adt/termvar"
	"orglang/go-engine/adt/uniqsym"
)

type API interface {
	Incept(uniqsym.ADT) (semtype.TypeRef, error)
	Create(DecSpec) (DecSnap, error)
	RetrieveSnap(semtype.TypeRef) (DecSnap, error)
	RetreiveRefs() ([]semtype.TypeRef, error)
}

type DecSpec struct {
	DescQN uniqsym.ADT
	// endpoint where process acts as a provider
	LiabVar termvar.VarSpec
	// endpoints where process acts as a client
	AssetVars []termvar.VarSpec
}

type DecRec struct {
	DescRef   semtype.TypeRef
	LiabVar   termvar.VarRec
	AssetVars []termvar.VarRec
}

// aka ExpDec or ExpDecDef without expression
type DecSnap struct {
	DescRef   semtype.TypeRef
	LiabVar   termvar.VarRec
	AssetVars []termvar.VarRec
}

type service struct {
	procDecs Repo
	descSems semtype.Repo
	operator db.Operator
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(procDecs Repo, descSems semtype.Repo, operator db.Operator, log *slog.Logger) *service {
	return &service{procDecs, descSems, operator, log}
}

func (s *service) Incept(procQN uniqsym.ADT) (_ semtype.TypeRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", procQN)
	s.log.Debug("starting inception...", qnAttr)
	newRef := semtype.NewRef()
	newDec := DecRec{DescRef: newRef}
	newDesc := semtype.SemRec{DescRef: newRef, DescQN: procQN, Kind: semtype.Proc}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descSems.AddRec(ds, newDesc)
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
		return semtype.TypeRef{}, err
	}
	s.log.Debug("inception succeed", qnAttr, slog.Any("ref", newDec.DescRef))
	return newDec.DescRef, nil
}

func (s *service) Create(spec DecSpec) (_ DecSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", spec.DescQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	typeQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars)+1)
	for _, spec := range spec.AssetVars {
		typeQNs = append(typeQNs, spec.TypeQN)
	}
	var typeRefs map[uniqsym.ADT]semtype.TypeRef
	getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		typeRefs, err = s.descSems.GetRefsByQNs(ds, append(typeQNs, spec.LiabVar.TypeQN))
		return err
	})
	if getErr != nil {
		return DecSnap{}, getErr
	}
	newLiabVar := termvar.VarRec{
		TypeRef: typeRefs[spec.LiabVar.TypeQN],
		ChnlPH:  spec.LiabVar.ChnlPH,
	}
	newAssetVars := make([]termvar.VarRec, 0, len(spec.AssetVars))
	for _, vs := range spec.AssetVars {
		newAssetVars = append(newAssetVars, termvar.VarRec{
			TypeRef: typeRefs[vs.TypeQN],
			ChnlPH:  vs.ChnlPH,
		})
	}
	newRef := semtype.NewRef()
	newDesc := semtype.SemRec{DescRef: newRef, DescQN: spec.DescQN, Kind: semtype.Proc}
	newDec := DecRec{DescRef: newRef, LiabVar: newLiabVar, AssetVars: newAssetVars}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descSems.AddRec(ds, newDesc)
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
	return DecSnap{DescRef: newRef, LiabVar: newLiabVar, AssetVars: newAssetVars}, nil
}

func (s *service) RetrieveSnap(ref semtype.TypeRef) (snap DecSnap, err error) {
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

func (s *service) RetreiveRefs() (refs []semtype.TypeRef, err error) {
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
