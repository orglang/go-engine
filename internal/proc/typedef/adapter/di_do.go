package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/proc/typedef/core"
	typeexp "orglang/go-engine/internal/proc/typeexp/core"
)

// Package registers process type-definition persistence and service.
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
		typeDefRepo, err := do.Invoke[core.Repo](i)
		if err != nil {
			return nil, err
		}
		typeExpRepo, err := do.Invoke[typeexp.Repo](i)
		if err != nil {
			return nil, err
		}
		descSemRepo, err := do.Invoke[descsem.Repo](i)
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
		return core.NewService(typeDefRepo, typeExpRepo, descSemRepo, transactor, logger), nil
	}),
)
