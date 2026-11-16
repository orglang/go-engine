package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/termdec/core"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

// Package registers process term-declaration persistence and service.
var Package = do.Package(
	do.Lazy[core.QueryBuilder](func(do.Injector) (core.QueryBuilder, error) {
		return newSQLBuilder(), nil
	}),
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
		builder, err := do.Invoke[core.QueryBuilder](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoPgx(builder, logger), nil
	}),
	do.Lazy[core.API](func(i do.Injector) (core.API, error) {
		termDecRepo, err := do.Invoke[core.Repo](i)
		if err != nil {
			return nil, err
		}
		typeDefRepo, err := do.Invoke[typedef.Repo](i)
		if err != nil {
			return nil, err
		}
		typeSemRepo, err := do.Invoke[typesem.Repo](i)
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
		return core.NewService(termDecRepo, typeDefRepo, typeSemRepo, transactor, logger), nil
	}),
)
