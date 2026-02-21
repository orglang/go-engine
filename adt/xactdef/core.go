package xactdef

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Create(DefSpec) (DefSnap, error)
}

type DefSpec struct {
	XactQN uniqsym.ADT
	XactES xactexp.ExpSpec
}

type DefRec struct {
	DescRef descsem.SemRef
	ExpVK   valkey.ADT
}

type DefSnap struct {
	DescRef descsem.SemRef
	DefSpec DefSpec
}

type service struct {
	xactDefs Repo
	xactExps xactexp.Repo
	descSems descsem.Repo
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
	descSems descsem.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	return &service{xactDefs, xactExps, descSems, operator, log}
}

func (s *service) Create(spec DefSpec) (_ DefSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", spec.XactQN)
	s.log.Debug("starting creation...", qnAttr, slog.Any("spec", spec))
	newExp, convertErr := xactexp.ConvertSpecToRec(spec.XactES)
	if convertErr != nil {
		return DefSnap{}, convertErr
	}
	newRef := descsem.NewRef()
	newBind := descsem.SemBind{DescQN: spec.XactQN, DescID: newRef.DescID}
	newDesc := descsem.SemRec{Ref: newRef, Bind: newBind, Kind: descsem.Xact}
	newDef := DefRec{DescRef: newRef, ExpVK: newExp.Key()}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descSems.InsertRec(ds, newDesc)
		if err != nil {
			return err
		}
		err = s.xactExps.InsertRec(ds, newExp, newRef)
		if err != nil {
			return err
		}
		return s.xactDefs.InsertRec(ds, newDef)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return DefSnap{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newRef))
	return DefSnap{
		DescRef: newRef,
		DefSpec: spec,
	}, nil
}

func errOptimisticUpdate(got revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
