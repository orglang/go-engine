package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	compconn "orglang/go-engine/internal/pool/compconn/core"
)

// Package registers the pool connection repository and its private query builder.
var Package = do.Package(
	do.Lazy[compconn.QueryBuilder](func(do.Injector) (compconn.QueryBuilder, error) {
		return newSQLBuilder(), nil
	}),
	do.Lazy[compconn.Repo](func(i do.Injector) (compconn.Repo, error) {
		qb, err := do.Invoke[compconn.QueryBuilder](i)
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
