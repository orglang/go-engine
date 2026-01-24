package pooldec

import (
	"log/slog"

	"orglang/go-runtime/lib/db"

	"orglang/go-runtime/adt/identity"
	"orglang/go-runtime/adt/procbind"
	"orglang/go-runtime/adt/uniqsym"
)

// Port
type API interface {
	Create(DecSpec) (DecRef, error)
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

type DecSpec struct {
	PoolQN uniqsym.ADT
	// Endpoints where pool acts as a provider for insiders
	InsiderProvisionBCs []procbind.BindSpec
	// Endpoints where pool acts as a client for insiders
	InsiderReceptionBCs []procbind.BindSpec
	// Endpoints where pool acts as a provider for outsiders
	OutsiderProvisionBCs []procbind.BindSpec
	// Endpoints where pool acts as a client for outsiders
	OutsiderReceptionBCs []procbind.BindSpec
}

type DecRef struct {
	DecID identity.ADT
}

type DecRec struct {
	DecID identity.ADT
}

type service struct {
	poolDecs repo
	operator db.Operator
	log      *slog.Logger
}

func (s *service) Create(spec DecSpec) (DecRef, error) {
	return DecRef{}, nil
}
