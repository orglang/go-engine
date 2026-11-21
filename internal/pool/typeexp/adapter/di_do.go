package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/pool/typeexp/core"
)

// Package registers the pool type-expression repository and its private query builder.
var Package = do.Package(
	do.Lazy[QueryBuilder](func(do.Injector) (QueryBuilder, error) {
		return newSQLBuilder(), nil
	}),
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
		qb, err := do.Invoke[QueryBuilder](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoPgx(qb, logger), nil
	}),
)
