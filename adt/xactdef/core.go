package xactdef

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/syndec"
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
	ExpID  valkey.ADT
}

type DefSnap struct {
	DefRef  DefRef
	DefSpec DefSpec
}

type service struct {
	xactDefs Repo
	xactExps xactexp.Repo
	synDecs  syndec.Repo
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
	synDecs syndec.Repo,
	operator db.Operator,
	l *slog.Logger,
) *service {
	return &service{xactDefs, xactExps, synDecs, operator, l}
}

func (s *service) Create(spec DefSpec) (_ DefSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("xactQN", spec.XactQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newSyn := syndec.DecRec{DecQN: spec.XactQN, DecID: identity.New(), DecRN: revnum.New()}
	newExp, err := xactexp.ConvertSpecToRec(spec.XactES)
	if err != nil {
		return DefSnap{}, err
	}
	newType := DefRec{
		DefRef: DefRef{ID: newSyn.DecID, RN: newSyn.DecRN},
		ExpID:  newExp.Key(),
	}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.synDecs.InsertRec(ds, newSyn)
		if err != nil {
			return err
		}
		err = s.xactExps.InsertRec(ds, newExp)
		if err != nil {
			return err
		}
		err = s.xactDefs.Insert(ds, newType)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("creation failed", qnAttr)
		return DefSnap{}, err
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("defRef", newType.DefRef))
	return DefSnap{
		DefRef:  newType.DefRef,
		DefSpec: spec,
	}, nil
}

func errOptimisticUpdate(got revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
