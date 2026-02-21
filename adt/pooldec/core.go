package pooldec

import (
	"context"
	"log/slog"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/descvar"
	"orglang/go-engine/adt/uniqsym"
)

type API interface {
	Create(DecSpec) (descsem.SemRef, error)
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

type DecSpec struct {
	DescQN     uniqsym.ADT
	ProviderVS descvar.VarSpec
	ClientVSes []descvar.VarSpec
}

type DecRec struct {
	DescRef    descsem.SemRef
	ProviderVR descvar.VarRec
	ClientVRs  []descvar.VarRec
}

func newService(
	poolDecs Repo,
	descSems descsem.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	return &service{poolDecs, descSems, operator, log}
}

type service struct {
	poolDecs Repo
	descSems descsem.Repo
	operator db.Operator
	log      *slog.Logger
}

func (s *service) Create(spec DecSpec) (_ descsem.SemRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", spec.DescQN)
	s.log.Debug("starting creation...", qnAttr, slog.Any("spec", spec))
	xactQNs := make([]uniqsym.ADT, 0, len(spec.ClientVSes)+1)
	for _, spec := range spec.ClientVSes {
		xactQNs = append(xactQNs, spec.DescQN)
	}
	var xactRefs map[uniqsym.ADT]descsem.SemRef
	selectErr := s.operator.Implicit(ctx, func(ds db.Source) error {
		xactRefs, err = s.descSems.SelectRefsByQNs(ds, append(xactQNs, spec.ProviderVS.DescQN))
		return err
	})
	if selectErr != nil {
		return descsem.SemRef{}, selectErr
	}
	providerVR := descvar.VarRec{
		ChnlPH: spec.ProviderVS.ChnlPH,
		DescID: xactRefs[spec.ProviderVS.DescQN].DescID,
	}
	clientVRs := make([]descvar.VarRec, 0, len(spec.ClientVSes))
	for _, vs := range spec.ClientVSes {
		clientVRs = append(clientVRs, descvar.VarRec{ChnlPH: vs.ChnlPH, DescID: xactRefs[vs.DescQN].DescID})
	}
	newRef := descsem.NewRef()
	newBind := descsem.SemBind{DescQN: spec.DescQN, DescID: newRef.DescID}
	newDesc := descsem.SemRec{Ref: newRef, Bind: newBind, Kind: descsem.Pool}
	newDec := DecRec{
		DescRef:    newRef,
		ProviderVR: providerVR,
		ClientVRs:  clientVRs,
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.descSems.InsertRec(ds, newDesc)
		if err != nil {
			return err
		}
		return s.poolDecs.InsertRec(ds, newDec)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return descsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newRef))
	return newRef, nil
}
