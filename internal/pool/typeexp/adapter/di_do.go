package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/pool/ent"
	"orglang/go-engine/internal/pool/typeexp/core"
	typedefcore "orglang/go-engine/internal/pool/typedef/core"
)

// Package registers the pool type-expression repository.
var Package = do.Package(
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
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
	do.Lazy[typedefcore.TypeExpRepo](func(i do.Injector) (typedefcore.TypeExpRepo, error) {
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil { return nil, err }
		return NewDaoEnt(client, logger), nil
	}),
)
