package termdec

import (
	"context"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/semterm"
	"orglang/go-engine/adt/termvar"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/pool/typedef"
)

type API interface {
	Create(DecSpec) (semterm.TermRef, error)
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

type DecSpec struct {
	TermQN    uniqsym.ADT
	LiabVar   termvar.VarSpec
	AssetVars []termvar.VarSpec
}

type DecRec struct {
	TermRef   semterm.TermRef
	LiabVar   termvar.VarRec
	AssetVars []termvar.VarRec
}

type DecSnap struct {
	TermRef semterm.TermRef
}

func newService(
	termDecs Repo,
	typeDefs typedef.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	return &service{termDecs, typeDefs, operator, log}
}

type service struct {
	termDecs Repo
	typeDefs typedef.Repo
	operator db.Operator
	log      *slog.Logger
}

func (s *service) Create(spec DecSpec) (_ semterm.TermRef, err error) {
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
		return semterm.TermRef{}, getErr
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
	newDec := DecRec{TermRef: semterm.NewRef(), LiabVar: newLiabVar, AssetVars: newAssetVars}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		return s.termDecs.AddRec(ds, newDec)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return semterm.TermRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newDec.TermRef))
	return newDec.TermRef, nil
}
