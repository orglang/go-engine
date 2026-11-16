package core

import (
	"log/slog"

	pooltypedef "orglang/go-engine/internal/pool/typedef/core"
	proctypedef "orglang/go-engine/internal/proc/typedef/core"
)

type API interface {
	Create(Spec) error
}

type Spec struct {
	Pools []pooltypedef.DefSpec
	Procs []proctypedef.DefSpec
}

type service struct {
	log *slog.Logger
}

// NewService constructs the program service.\nfunc NewService(log *slog.Logger) API {
	return &service{log: log}
}

func (s *service) Create(spec Spec) error {
	return nil
}
