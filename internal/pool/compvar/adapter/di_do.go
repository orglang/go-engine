package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	compvar "orglang/go-engine/internal/pool/compvar/core"
)

// Package registers the pool component-variable repository and its private query builder.
var Package = do.Package(
	do.Lazy[compvar.QueryBuilder](func(do.Injector) (compvar.QueryBuilder, error) {
		return newSQLBuilder(), nil
	}),
	do.Lazy[compvar.Repo](func(i do.Injector) (compvar.Repo, error) {
		qb, err := do.Invoke[compvar.QueryBuilder](i)
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
