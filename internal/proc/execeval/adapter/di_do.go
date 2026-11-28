package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	execcomm "orglang/go-engine/internal/proc/execcomm/core"
	"orglang/go-engine/internal/proc/execeval/core"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

// Package registers the process computation-execution repository, service, and HTTP controller.
var Package = do.Package(
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return newDaoEnt(client, logger), nil
	}),
	do.Lazy[core.API](func(i do.Injector) (core.API, error) {
		evalRepo, err := do.Invoke[core.Repo](i)
		if err != nil {
			return nil, err
		}
		commRepo, err := do.Invoke[execcomm.Repo](i)
		if err != nil {
			return nil, err
		}
		termDecRepo, err := do.Invoke[termdec.Repo](i)
		if err != nil {
			return nil, err
		}
		typeDefRepo, err := do.Invoke[typedef.Repo](i)
		if err != nil {
			return nil, err
		}
		typeExprrRepo, err := do.Invoke[typedef.TypeExprRepo](i)
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
		return core.NewService(evalRepo, commRepo, termDecRepo, typeDefRepo, typeExprrRepo, transactor, logger), nil
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
