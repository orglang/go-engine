package pooldecl

import (
	"log/slog"

	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
)

// Port
type API interface {
	Create(PoolSpec) (PoolRef, error)
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

type PoolSpec struct {
	PoolNS qualsym.ADT
	PoolSN qualsym.ADT
	// endpoint where pool acts as a provider for insiders
	InsiderProvisionEP ChnlSpec
	// endpoint where pool acts as a client for insiders
	InsiderReceptionEP ChnlSpec
	// endpoint where pool acts as a provider for outsiders
	OutsiderProvisionEP ChnlSpec
	// endpoints where pool acts as a client for outsiders
	OutsiderReceptionEPs []ChnlSpec
}

type ChnlSpec struct {
	CommPH qualsym.ADT // may be blank
	TypeQN qualsym.ADT
}

type PoolRef struct {
	DecID identity.ADT
}

type poolRec struct {
	DecID identity.ADT
}

type service struct {
	sigs     repo
	operator sd.Operator
	log      *slog.Logger
}

func (s *service) Create(spec PoolSpec) (PoolRef, error) {
	return PoolRef{}, nil
}
