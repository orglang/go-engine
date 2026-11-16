package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/typedef/core"
	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

// Package registers the type-definition service, persistence, and HTTP adapter.
var Package = do.Package(
	do.Lazy[core.QueryBuilder](func(do.Injector) (core.QueryBuilder, error) {
		return newSQLBuilder(), nil
	}),
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
		qb, err := do.Invoke[core.QueryBuilder](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoPgx(qb, logger), nil
	}),
	do.Lazy[core.API](func(i do.Injector) (core.API, error) {
		repo, err := do.Invoke[core.Repo](i)
		if err != nil {
			return nil, err
		}
		typeExpRepo, err := do.Invoke[typeexp.Repo](i)
		if err != nil {
			return nil, err
		}
		descRepo, err := do.Invoke[descsem.Repo](i)
		if err != nil {
			return nil, err
		}
		transactor, err := do.Invoke[tx.Transactor](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return core.NewService(repo, typeExpRepo, descRepo, transactor, logger), nil
	}),
	do.Lazy[Controller](func(i do.Injector) (Controller, error) {
		api, err := do.Invoke[core.API](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return newControllerEcho(api, logger), nil
	}),
)
