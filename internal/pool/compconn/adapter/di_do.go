package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/pool/ent"
	compconn "orglang/go-engine/internal/pool/compconn/core"
)

// Package registers the pool connection repository.
var Package = do.Package(
	do.Lazy[compconn.Repo](func(i do.Injector) (compconn.Repo, error) {
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
