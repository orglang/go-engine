package adapter

import (
	"context"
	"log/slog"

	"github.com/alitto/pond/v2"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/lib/tx"
	compvar "orglang/go-engine/internal/pool/compvar/core"
	"orglang/go-engine/internal/pool/ent"
	execconn "orglang/go-engine/internal/pool/execconn/core"
	execeval "orglang/go-engine/internal/pool/execeval/core"
	termdec "orglang/go-engine/internal/pool/termdec/core"
	typedef "orglang/go-engine/internal/pool/typedef/core"
	procexeceval "orglang/go-engine/internal/proc/execeval/core"
)

// Package registers the computation-execution repository, service, broker and
// inbound controller. Broker subscription and route registration are explicit
// composition steps.
var Package = do.Package(
	do.Lazy[execeval.Repo](func(i do.Injector) (execeval.Repo, error) {
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewDaoEnt(client, logger), nil
	}),
	do.Lazy[execeval.Broker](func(i do.Injector) (execeval.Broker, error) {
		pool, err := do.Invoke[pond.Pool](i)
		if err != nil {
			return nil, err
		}
		taskCtx, err := do.Invoke[context.Context](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return newPondBroker(taskCtx, pool, logger), nil
	}),
	do.Lazy[execeval.API](func(i do.Injector) (execeval.API, error) {
		evalRepo, err := do.Invoke[execeval.Repo](i)
		if err != nil {
			return nil, err
		}
		broker, err := do.Invoke[execeval.Broker](i)
		if err != nil {
			return nil, err
		}
		varRepo, err := do.Invoke[compvar.Repo](i)
		if err != nil {
			return nil, err
		}
		connRepo, err := do.Invoke[execconn.Repo](i)
		if err != nil {
			return nil, err
		}
		typeRepo, err := do.Invoke[typedef.Repo](i)
		if err != nil {
			return nil, err
		}
		procRepo, err := do.Invoke[procexeceval.Repo](i)
		if err != nil {
			return nil, err
		}
		termDecRepo, err := do.Invoke[termdec.Repo](i)
		if err != nil {
			return nil, err
		}
		compRepo, err := do.Invoke[comproot.Repo](i)
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
		return execeval.NewService(execeval.ServiceDeps{
			ExecEvalRepo:   evalRepo,
			ExecEvalBroker: broker,
			CompVarRepo:    varRepo,
			ExecConnRepo:   connRepo,
			TypeDefRepo:    typeRepo,
			ProcEvalRepo:   procRepo,
			TermDecRepo:    termDecRepo,
			CompRootRepo:   compRepo,
			Transactor:     transactor,
			Logger:         logger,
		}), nil
	}),
	do.Lazy[Controller](func(i do.Injector) (Controller, error) {
		api, err := do.Invoke[execeval.API](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return newControllerEcho(api, logger), nil
	}),
	do.Lazy[comproot.Repo](func(i do.Injector) (comproot.Repo, error) {
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewCompRootDaoEnt(client, logger), nil
	}),
)

// SubscribeBroker connects the asynchronous broker to the resolved service.
// The application composition root must call this after registering Package.
func SubscribeBroker(i do.Injector) error {
	broker, err := do.Invoke[execeval.Broker](i)
	if err != nil {
		return err
	}
	api, err := do.Invoke[execeval.API](i)
	if err != nil {
		return err
	}
	broker.Subscribe(api)
	return nil
}
