package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	webapi "orglang/go-engine/internal/app/web/api"
	"orglang/go-engine/internal/lib/tx"
	compconn "orglang/go-engine/internal/proc/compconn/core"
	"orglang/go-engine/internal/proc/compexec/core"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	typedef "orglang/go-engine/internal/proc/typedef/core"
	typeexp "orglang/go-engine/internal/proc/typeexp/core"
)

// Package registers the process computation-execution repository, service, and HTTP controller.
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
		return newDaoPgx(qb, logger), nil
	}),
	do.Lazy[core.API](func(i do.Injector) (core.API, error) {
		execRepo, err := do.Invoke[core.Repo](i)
		if err != nil {
			return nil, err
		}
		connRepo, err := do.Invoke[compconn.Repo](i)
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
		typeExpRepo, err := do.Invoke[typeexp.Repo](i)
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
		return core.NewService(execRepo, connRepo, termDecRepo, typeDefRepo, typeExpRepo, transactor, logger), nil
	}),
	do.Lazy[GeneratedController](func(i do.Injector) (GeneratedController, error) {
		api, err := do.Invoke[core.API](i)
		if err != nil { return nil, err }
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil { return nil, err }
		return newControllerEcho(api, logger), nil
	}),
	do.Lazy[Controller](func(i do.Injector) (Controller, error) {
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
