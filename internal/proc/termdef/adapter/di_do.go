package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/descsem"
	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"
	webapi "orglang/go-engine/internal/app/web/api"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	"orglang/go-engine/internal/proc/termdef/core"
	termexp "orglang/go-engine/internal/proc/termexp/core"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

// Package registers the unified process term lifecycle service and its HTTP adapters.
var Package = do.Package(
	do.LazyNamed[descsem.Repo]("proc-descsem", func(i do.Injector) (descsem.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return descsemadapter.NewDaoPgx("proc_desc_binds ")(logger), nil
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
	do.Lazy[core.DefRepo](func(i do.Injector) (core.DefRepo, error) {
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
		termDefRepo, err := do.Invoke[core.Repo](i)
		if err != nil {
			return nil, err
		}
		defRepo, err := do.Invoke[core.DefRepo](i)
		if err != nil {
			return nil, err
		}
		termExpRepo, err := do.Invoke[termexp.Repo](i)
		if err != nil {
			return nil, err
		}
		descRootRepo, err := do.InvokeNamed[descroot.Repo](i, "proc-descroot")
		if err != nil {
			return nil, err
		}
		descSemRepo, err := do.InvokeNamed[descsem.Repo](i, "proc-descsem")
		if err != nil {
			return nil, err
		}
		typeDefRepo, err := do.Invoke[typedef.Repo](i)
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
		return core.NewService(termDefRepo, defRepo, termExpRepo, descRootRepo, descSemRepo, typeDefRepo, transactor, logger), nil
	}),
	do.Lazy[webapi.ServerInterface](func(i do.Injector) (webapi.ServerInterface, error) {
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
