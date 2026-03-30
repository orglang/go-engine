package pooldec

import (
	"context"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/descvar"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/xactdef"
)

type API interface {
	Create(DecSpec) (descsem.SemRef, error)
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

type DecSpec struct {
	DescQN    uniqsym.ADT
	LiabVar   descvar.VarSpec
	AssetVars []descvar.VarSpec
}

type DecRec struct {
	DescRef   descsem.SemRef
	LiabVar   descvar.VarRec
	AssetVars []descvar.VarRec
}

type DecSnap struct {
	DescRef descsem.SemRef
}

func newService(
	poolDecs Repo,
	descSems descsem.Repo,
	xactDefs xactdef.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	return &service{poolDecs, descSems, xactDefs, operator, log}
}

type service struct {
	poolDecs Repo
	descSems descsem.Repo
	xactDefs xactdef.Repo
	operator db.Operator
	log      *slog.Logger
}

func (s *service) Create(spec DecSpec) (_ descsem.SemRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", spec.DescQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars)+1)
	for _, varSpec := range spec.AssetVars {
		assetQNs = append(assetQNs, varSpec.DescQN)
	}
	var xactDefs map[uniqsym.ADT]xactdef.DefRec
	getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		xactDefs, err = s.xactDefs.GetRecsByQNs(ds, append(assetQNs, spec.LiabVar.DescQN))
		return err
	})
	if getErr != nil {
		return descsem.SemRef{}, getErr
	}
	newLiabVar := descvar.VarRec{
		DescRef: xactDefs[spec.LiabVar.DescQN].DescRef,
		ChnlPH:  spec.LiabVar.ChnlPH,
		ExpVK:   xactDefs[spec.LiabVar.DescQN].ExpVK,
	}
	newAssetVars := make([]descvar.VarRec, 0, len(spec.AssetVars))
	for _, varSpec := range spec.AssetVars {
		newAssetVars = append(newAssetVars, descvar.VarRec{
			DescRef: xactDefs[varSpec.DescQN].DescRef,
			ChnlPH:  varSpec.ChnlPH,
			ExpVK:   xactDefs[varSpec.DescQN].ExpVK,
		})
	}
	newRef := descsem.NewRef()
	newDesc := descsem.SemRec{DescRef: newRef, DescQN: spec.DescQN, Kind: descsem.Pool}
	newDec := DecRec{DescRef: newRef, LiabVar: newLiabVar, AssetVars: newAssetVars}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descSems.AddRec(ds, newDesc)
		if err != nil {
			return err
		}
		return s.poolDecs.AddRec(ds, newDec)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return descsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newRef))
	return newRef, nil
}
