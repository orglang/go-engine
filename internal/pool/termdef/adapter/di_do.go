package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/termdef/core"
	typedef "orglang/go-engine/internal/pool/typedef/core"
)

// Package registers the term-definition service, persistence, and HTTP adapter.
var Package = do.Package(
	do.Lazy[queryBuilder](func(do.Injector) (queryBuilder, error) {
		return newSQLBuilder(), nil
	}),
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
		qb, err := do.Invoke[queryBuilder](i)
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
		typeDefRepo, err := do.Invoke[typedef.Repo](i)
		if err != nil {
			return nil, err
		}
		descRepo, err := do.InvokeNamed[descsem.Repo](i, "pool-descsem")
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
		return core.NewService(repo, typeDefRepo, descRepo, transactor, logger), nil
	}),
	do.Lazy[GeneratedController](func(i do.Injector) (GeneratedController, error) {
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
