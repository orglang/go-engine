package dec

import (
	"log/slog"

	"smecalculus/rolevod/lib/data"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type PoolSpec struct {
	X      ChnlSpec // export
	PoolNS sym.ADT
	PoolSN sym.ADT    // label
	Ys     []ChnlSpec // imports
}

type ChnlSpec struct {
	CommPH sym.ADT // may be blank
	TypeQN sym.ADT
}

type PoolRef struct {
	DecID id.ADT
}

type poolRec struct {
	DecID id.ADT
}

type API interface {
	Create(PoolSpec) (PoolRef, error)
}

// for compilation purposes
func newAPI() API {
	return &service{}
}

type service struct {
	sigs     repo
	operator data.Operator
	log      *slog.Logger
}

func (s *service) Create(spec PoolSpec) (PoolRef, error) {
	return PoolRef{}, nil
}

// Port
type repo interface {
	Insert(data.Source, poolRec) error
}
