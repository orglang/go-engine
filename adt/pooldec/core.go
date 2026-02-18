package pooldec

import (
	"context"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descbind"
	"orglang/go-engine/adt/descexec"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/poolbind"
	"orglang/go-engine/adt/uniqsym"
)

type API interface {
	Create(DecSpec) (descexec.ExecRef, error)
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

type DecSpec struct {
	PoolQN     uniqsym.ADT
	ProviderBS poolbind.BindSpec
	ClientBSes []poolbind.BindSpec
}

type DecRec struct {
	PoolID     identity.ADT
	ProviderBR poolbind.BindRec
	ClientBRs  []poolbind.BindRec
}

func newService(
	poolDecs repo,
	descExecs descexec.Repo,
	descBinds descbind.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	return &service{poolDecs, descExecs, descBinds, operator, log}
}

type service struct {
	poolDecs  repo
	descExecs descexec.Repo
	descBinds descbind.Repo
	operator  db.Operator
	log       *slog.Logger
}

func (s *service) Create(spec DecSpec) (_ descexec.ExecRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("poolQN", spec.PoolQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	xactQNs := make([]uniqsym.ADT, 0, len(spec.ClientBSes)+1)
	for _, bs := range spec.ClientBSes {
		xactQNs = append(xactQNs, bs.XactQN)
	}
	var xactDefs map[uniqsym.ADT]descexec.ExecRef
	selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		xactDefs, err = s.descExecs.SelectRefsByQNs(ds, append(xactQNs, spec.ProviderBS.XactQN))
		return err
	})
	if selectErr != nil {
		return descexec.ExecRef{}, selectErr
	}
	providerBR := poolbind.BindRec{
		ChnlPH: spec.ProviderBS.ChnlPH,
		DefID:  xactDefs[spec.ProviderBS.XactQN].DescID,
	}
	clientBRs := make([]poolbind.BindRec, 0, len(spec.ClientBSes))
	for _, bs := range spec.ClientBSes {
		clientBRs = append(clientBRs, poolbind.BindRec{ChnlPH: bs.ChnlPH, DefID: xactDefs[bs.XactQN].DescID})
	}
	newExec := descexec.ExecRec{Ref: descexec.NewExecRef(), Kind: descexec.Pool}
	newBind := descbind.BindRec{DescQN: spec.PoolQN, DescID: newExec.Ref.DescID}
	newDec := DecRec{
		PoolID:     newExec.Ref.DescID,
		ProviderBR: providerBR,
		ClientBRs:  clientBRs,
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descBinds.InsertRec(ds, newBind)
		if err != nil {
			return err
		}
		err = s.descExecs.InsertRec(ds, newExec)
		if err != nil {
			return err
		}
		return s.poolDecs.InsertRec(ds, newDec)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return descexec.ExecRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("poolRef", newExec.Ref))
	return newExec.Ref, nil
}
