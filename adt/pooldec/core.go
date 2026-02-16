package pooldec

import (
	"context"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/poolbind"
	"orglang/go-engine/adt/synonym"
	"orglang/go-engine/adt/uniqref"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
	"orglang/go-engine/adt/xactdef"
)

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

func newService(
	poolDecs repo,
	xactDefs xactdef.Repo,
	synonyms synonym.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	return &service{poolDecs, xactDefs, synonyms, operator, log}
}

type service struct {
	poolDecs repo
	xactDefs xactdef.Repo
	synonyms synonym.Repo
	operator db.Operator
	log      *slog.Logger
}

func (s *service) Create(spec DecSpec) (_ DecRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("poolQN", spec.PoolQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	xactQNs := make([]uniqsym.ADT, 0, len(spec.ClientBSes)+1)
	for _, bs := range spec.ClientBSes {
		xactQNs = append(xactQNs, bs.XactQN)
	}
	var xactDefs map[uniqsym.ADT]xactdef.DefRef
	selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		xactDefs, err = s.xactDefs.SelectRefsByQNs(ds, append(xactQNs, spec.ProviderBS.XactQN))
		return err
	})
	if selectErr != nil {
		return DecRef{}, selectErr
	}
	providerBR := poolbind.BindRec{
		ChnlPH: spec.ProviderBS.ChnlPH,
		DefID:  xactDefs[spec.ProviderBS.XactQN].ID,
	}
	clientBRs := make([]poolbind.BindRec, 0, len(spec.ClientBSes))
	for _, bs := range spec.ClientBSes {
		clientBRs = append(clientBRs, poolbind.BindRec{ChnlPH: bs.ChnlPH, DefID: xactDefs[bs.XactQN].ID})
	}
	synVK, keyErr := spec.PoolQN.Key()
	if keyErr != nil {
		return DecRef{}, keyErr
	}
	newSyn := synonym.Rec{SynQN: spec.PoolQN, SynVK: synVK}
	newDec := DecRec{
		DecRef:     uniqref.New(),
		SynVK:      synVK,
		ProviderBR: providerBR,
		ClientBRs:  clientBRs,
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.synonyms.InsertRec(ds, newSyn)
		if err != nil {
			return err
		}
		return s.poolDecs.InsertRec(ds, newDec)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return DecRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("decRef", newDec.DecRef))
	return newDec.DecRef, nil
}
