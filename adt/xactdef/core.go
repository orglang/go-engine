package xactdef

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/seqnum"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Create(DefSpec) (descsem.SemRef, error)
}

type DefSpec struct {
	XactQN  uniqsym.ADT
	XactExp xactexp.ExpSpec
}

type DefRec struct {
	DescRef descsem.SemRef
	ExpVK   valkey.ADT
}

type DefSnap struct {
	DescRef descsem.SemRef
	XactExp xactexp.ExpSpec
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

func (s *service) Create(spec DefSpec) (_ descsem.SemRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", spec.XactQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newExp, convErr := xactexp.ConvertSpecToRec(spec.XactExp)
	if convErr != nil {
		return descsem.SemRef{}, convErr
	}
	newRef := descsem.NewRef()
	newDesc := descsem.SemRec{DescRef: newRef, DescQN: spec.XactQN, Kind: descsem.Xact}
	newDef := DefRec{DescRef: newRef, ExpVK: newExp.Key()}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descSems.AddRec(ds, newDesc)
		if err != nil {
			return err
		}
		err = s.xactExps.AddRec(ds, newExp, newRef)
		if err != nil {
			return err
		}
		return s.xactDefs.InsertRec(ds, newDef)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return descsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("xact", newRef))
	return newRef, nil
}

func errOptimisticUpdate(got seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}

func ErrMissingInEnv(want valkey.ADT) error {
	return fmt.Errorf("exp missing in env: %v", want)
}
