package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/progroot"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	"orglang/go-engine/internal/proc/typedef/core"
)

// Package registers process type-definition persistence, service, and HTTP adapters.
var Package = do.Package(
	do.Lazy[core.TypeExpRepo](func(i do.Injector) (core.TypeExpRepo, error) {
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewExpDaoEnt(client, logger), nil
	}),
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		return NewDaoEnt(client, logger), nil
	}),
	do.Lazy[core.API](func(i do.Injector) (core.API, error) {
		typeDefRepo, err := do.Invoke[core.Repo](i)
		if err != nil {
			return nil, err
		}
		progRootRepo, err := do.InvokeNamed[progroot.Repo](i, "proc-progroot")
		if err != nil {
			return nil, err
		}
		typeExpRepo, err := do.Invoke[core.TypeExpRepo](i)
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
		return core.NewService(typeDefRepo, progRootRepo, typeExpRepo, transactor, logger), nil
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
