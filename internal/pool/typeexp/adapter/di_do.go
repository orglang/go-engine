package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	core "orglang/go-engine/internal/pool/typeexp/core"
)

// Package registers the pool type-expression repository and its private query builder.
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
)
