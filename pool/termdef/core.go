package termdef

import (
	"context"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/termsem"
	"orglang/go-engine/adt/termvar"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/pool/typedef"
)

type API interface {
	Create(DefSpec) (termsem.SemRef, error)
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

type DefSpec struct {
	TermQN    uniqsym.ADT
	LiabVar   termvar.VarSpec
	AssetVars []termvar.VarSpec
}

type DefRec struct {
	TermRef   termsem.SemRef
	LiabVar   termvar.VarRec
	AssetVars []termvar.VarRec
}

type DefSnap struct {
	TermRef termsem.SemRef
}

func newService(
	termDefs Repo,
	typeDefs typedef.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	return &service{termDefs, typeDefs, operator, log}
}

type service struct {
	termDefs Repo
	typeDefs typedef.Repo
	operator db.Operator
	log      *slog.Logger
}

func (s *service) Create(spec DefSpec) (_ termsem.SemRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", spec.TermQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars)+1)
	for _, varSpec := range spec.AssetVars {
		assetQNs = append(assetQNs, varSpec.TypeQN)
	}
	var typeDefs map[uniqsym.ADT]typedef.DefRec
	getErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		typeDefs, err = s.typeDefs.GetRecsByQNs(ds, append(assetQNs, spec.LiabVar.TypeQN))
		return err
	})
	if getErr != nil {
		return termsem.SemRef{}, getErr
	}
	newLiabVar := termvar.VarRec{
		TypeRef: typeDefs[spec.LiabVar.TypeQN].TypeRef,
		ChnlPH:  spec.LiabVar.ChnlPH,
		ExpVK:   typeDefs[spec.LiabVar.TypeQN].ExpVK,
	}
	newAssetVars := make([]termvar.VarRec, 0, len(spec.AssetVars))
	for _, varSpec := range spec.AssetVars {
		newAssetVars = append(newAssetVars, termvar.VarRec{
			TypeRef: typeDefs[varSpec.TypeQN].TypeRef,
			ChnlPH:  varSpec.ChnlPH,
			ExpVK:   typeDefs[varSpec.TypeQN].ExpVK,
		})
	}
	newDef := DefRec{TermRef: termsem.New(), LiabVar: newLiabVar, AssetVars: newAssetVars}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		return s.termDefs.AddRec(ds, newDef)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return termsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newDef.TermRef))
	return newDef.TermRef, nil
}
