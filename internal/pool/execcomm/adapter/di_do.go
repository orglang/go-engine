package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/pool/ent"
	execcomm "orglang/go-engine/internal/pool/execcomm/core"
)

// Package registers the pool communication repository.
var Package = do.Package(
	do.Lazy[execcomm.Repo](func(i do.Injector) (execcomm.Repo, error) {
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoEnt(client, logger), nil
	}),
)
