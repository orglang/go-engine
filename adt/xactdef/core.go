package xactdef

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/synonym"
	"orglang/go-engine/adt/uniqref"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Create(DefSpec) (DefSnap, error)
}

type DefRef = uniqref.ADT

type DefSpec struct {
	XactQN uniqsym.ADT
	XactES xactexp.ExpSpec
}

type DefRec struct {
	DefRef DefRef
	SynVK  valkey.ADT
	ExpVK  valkey.ADT
}

type DefSnap struct {
	DefRef  DefRef
	DefSpec DefSpec
}

type service struct {
	xactDefs Repo
	xactExps xactexp.Repo
	synonyms synonym.Repo
	operator db.Operator
	log      *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	xactDefs Repo,
	xactExps xactexp.Repo,
	synonyms synonym.Repo,
	operator db.Operator,
	l *slog.Logger,
) *service {
	return &service{xactDefs, xactExps, synonyms, operator, l}
}

func (s *service) Create(spec DefSpec) (_ DefSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("xactQN", spec.XactQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	synVK, keyErr := spec.XactQN.Key()
	if keyErr != nil {
		return DefSnap{}, keyErr
	}
	newSyn := synonym.Rec{SynQN: spec.XactQN, SynVK: synVK}
	newExp, convertErr := xactexp.ConvertSpecToRec(spec.XactES)
	if convertErr != nil {
		return DefSnap{}, convertErr
	}
	newDef := DefRec{
		DefRef: uniqref.New(),
		SynVK:  synVK,
		ExpVK:  newExp.Key(),
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.synonyms.InsertRec(ds, newSyn)
		if err != nil {
			return err
		}
		err = s.xactExps.InsertRec(ds, newExp, newDef.DefRef)
		if err != nil {
			return err
		}
		return s.xactDefs.InsertRec(ds, newDef)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return DefSnap{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("defRef", newDef.DefRef))
	return DefSnap{
		DefRef:  newDef.DefRef,
		DefSpec: spec,
	}, nil
}

func errOptimisticUpdate(got revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
