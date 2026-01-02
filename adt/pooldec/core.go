package pooldec

import (
	"log/slog"

	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
	"orglang/orglang/adt/termctx"
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
	InsiderProvisionEP termctx.BindClaim
	// endpoint where pool acts as a client for insiders
	InsiderReceptionEP termctx.BindClaim
	// endpoint where pool acts as a provider for outsiders
	OutsiderProvisionEP termctx.BindClaim
	// endpoints where pool acts as a client for outsiders
	OutsiderReceptionEPs []termctx.BindClaim
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
