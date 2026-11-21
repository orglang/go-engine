package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/typedef/core"
	typeexp "orglang/go-engine/internal/proc/typeexp/core"
)

// Package registers process type-definition persistence, service, and HTTP adapters.
var Package = do.Package(
	do.Lazy[QueryBuilder](func(do.Injector) (QueryBuilder, error) {
		return newSQLBuilder(), nil
	}),
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
		builder, err := do.Invoke[QueryBuilder](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoPgx(builder, logger), nil
	}),
	do.Lazy[core.DescSemRepo](func(i do.Injector) (core.DescSemRepo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return descsemadapter.NewDaoPgx("proc_desc_binds ")(logger), nil
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
		descSemRepo, err := do.Invoke[core.DescSemRepo](i)
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
	do.Lazy[Presenter](func(i do.Injector) (Presenter, error) {
		api, err := do.Invoke[core.API](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		renderer, err := newRendererStdlib(logger)
		if err != nil {
			return nil, err
		}
		return newEchoPresenter(api, renderer, logger), nil
	}),
)
