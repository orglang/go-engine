package pooldec

import (
	"context"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/poolbind"
	"orglang/go-engine/adt/uniqref"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
	"orglang/go-engine/adt/xactdef"
)

// Port
type API interface {
	Create(DecSpec) (DecRef, error)
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

type DecRef = uniqref.ADT

type DecRec struct {
	DecRef     DecRef
	SynVK      valkey.ADT
	ProviderBR poolbind.BindRec
	ClientBRs  []poolbind.BindRec
}

func newService(poolDecs repo, xactDefs xactdef.Repo, operator db.Operator, log *slog.Logger) *service {
	return &service{poolDecs, xactDefs, operator, log}
}

type service struct {
	poolDecs repo
	xactDefs xactdef.Repo
	operator db.Operator
	log      *slog.Logger
}

func (s *service) Create(spec DecSpec) (_ DecRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("poolQN", spec.PoolQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	xactQNs := make([]uniqsym.ADT, len(spec.ClientBSes)+1)
	for _, bs := range spec.ClientBSes {
		xactQNs = append(xactQNs, bs.XactQN)
	}
	var xactDefs map[uniqsym.ADT]xactdef.DefRef
	err = s.operator.Implicit(ctx, func(ds db.Source) error {
		xactDefs, err = s.xactDefs.SelectRefsByQNs(ds, append(xactQNs, spec.ProviderBS.XactQN))
		if err != nil {
			return err
		}
		return nil
	})
	providerBR := poolbind.BindRec{
		ChnlPH: spec.ProviderBS.ChnlPH,
		DefID:  xactDefs[spec.ProviderBS.XactQN].ID,
	}
	clientBRs := make([]poolbind.BindRec, len(spec.ClientBSes))
	for _, bs := range spec.ClientBSes {
		clientBRs = append(clientBRs, poolbind.BindRec{ChnlPH: bs.ChnlPH, DefID: xactDefs[bs.XactQN].ID})
	}
	synVK, err := spec.PoolQN.Key()
	if err != nil {
		return DecRef{}, err
	}
	newDec := DecRec{
		DecRef:     uniqref.New(),
		SynVK:      synVK,
		ProviderBR: providerBR,
		ClientBRs:  clientBRs,
	}
	err = s.operator.Explicit(ctx, func(ds db.Source) error {
		err := s.poolDecs.InsertRec(ds, newDec)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Error("creation failed", qnAttr)
		return DecRef{}, err
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("decRef", newDec.DecRef))
	return newDec.DecRef, nil
}
