package xactdef

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descbind"
	"orglang/go-engine/adt/descexec"
	"orglang/go-engine/adt/identity"
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
	XactID identity.ADT
	ExpVK  valkey.ADT
}

type DefSnap struct {
	DescRef descexec.ExecRef
	DefSpec DefSpec
}

type service struct {
	xactDefs  Repo
	xactExps  xactexp.Repo
	descExecs descexec.Repo
	descBinds descbind.Repo
	operator  db.Operator
	log       *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	xactDefs Repo,
	xactExps xactexp.Repo,
	descExecs descexec.Repo,
	descBinds descbind.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	return &service{xactDefs, xactExps, descExecs, descBinds, operator, log}
}

func (s *service) Create(spec DefSpec) (_ DefSnap, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("xactQN", spec.XactQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newExp, convertErr := xactexp.ConvertSpecToRec(spec.XactES)
	if convertErr != nil {
		return DefSnap{}, convertErr
	}
	newExec := descexec.ExecRec{Ref: descexec.NewExecRef(), Kind: descexec.Xact}
	newBind := descbind.BindRec{DescQN: spec.XactQN, DescID: newExec.Ref.DescID}
	newDef := DefRec{XactID: newExec.Ref.DescID, ExpVK: newExp.Key()}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descBinds.InsertRec(ds, newBind)
		if err != nil {
			return err
		}
		err = s.descExecs.InsertRec(ds, newExec)
		if err != nil {
			return err
		}
		err = s.xactExps.InsertRec(ds, newExp, newExec.Ref)
		if err != nil {
			return err
		}
		return s.xactDefs.InsertRec(ds, newDef)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return DefSnap{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("xactRef", newExec.Ref))
	return DefSnap{
		DescRef: newExec.Ref,
		DefSpec: spec,
	}, nil
}

func errOptimisticUpdate(got revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
