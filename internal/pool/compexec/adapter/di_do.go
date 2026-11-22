package adapter

import (
	"context"
	"log/slog"

	"github.com/alitto/pond/v2"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/implsem"
	implsemadapter "orglang/go-engine/internal/adt/implsem/adapter"
	"orglang/go-engine/internal/lib/tx"
	compconn "orglang/go-engine/internal/pool/compconn/core"
	compexec "orglang/go-engine/internal/pool/compexec/core"
	compvar "orglang/go-engine/internal/pool/compvar/core"
	"orglang/go-engine/internal/pool/ent"
	termdef "orglang/go-engine/internal/pool/termdef/core"
	typeexp "orglang/go-engine/internal/pool/typeexp/core"
	proccompexec "orglang/go-engine/internal/proc/compexec/core"
)

// Package registers the computation-execution repository, service, broker and
// inbound controller. Broker subscription and route registration are explicit
// composition steps.
var Package = do.Package(
	do.Lazy[compexec.Repo](func(i do.Injector) (compexec.Repo, error) {
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
	do.Lazy[compexec.Broker](func(i do.Injector) (compexec.Broker, error) {
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
	do.Lazy[compexec.API](func(i do.Injector) (compexec.API, error) {
		execRepo, err := do.Invoke[compexec.Repo](i)
		if err != nil {
			return nil, err
		}
		broker, err := do.Invoke[compexec.Broker](i)
		if err != nil {
			return nil, err
		}
		varRepo, err := do.Invoke[compvar.Repo](i)
		if err != nil {
			return nil, err
		}
		connRepo, err := do.Invoke[compconn.Repo](i)
		if err != nil {
			return nil, err
		}
		typeRepo, err := do.Invoke[typeexp.Repo](i)
		if err != nil {
			return nil, err
		}
		procRepo, err := do.Invoke[proccompexec.Repo](i)
		if err != nil {
			return nil, err
		}
		termRepo, err := do.Invoke[termdef.Repo](i)
		if err != nil {
			return nil, err
		}
		implRepo, err := do.Invoke[implsem.Repo](i)
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
		return compexec.NewService(compexec.ServiceDeps{
			CompExecRepo:   execRepo,
			CompExecBroker: broker,
			CompVarRepo:    varRepo,
			CompConnRepo:   connRepo,
			TypeExpRepo:    typeRepo,
			ProcExecRepo:   procRepo,
			TermDefRepo:    termRepo,
			ImplSemRepo:    implRepo,
			CompSemRepo:    compRepo,
			Transactor:     transactor,
			Logger:         logger,
		}), nil
	}),
	do.Lazy[Controller](func(i do.Injector) (Controller, error) {
		api, err := do.Invoke[compexec.API](i)
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
		return NewCompSemDaoEnt(client, logger), nil
	}),
	do.Lazy[implsem.Repo](func(i do.Injector) (implsem.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return implsemadapter.NewDaoPgx("pool_impl_binds")(logger), nil
	}),
)

// SubscribeBroker connects the asynchronous broker to the resolved service.
// The application composition root must call this after registering Package.
func SubscribeBroker(i do.Injector) error {
	broker, err := do.Invoke[compexec.Broker](i)
	if err != nil {
		return err
	}
	api, err := do.Invoke[compexec.API](i)
	if err != nil {
		return err
	}
	broker.Subscribe(api)
	return nil
}
