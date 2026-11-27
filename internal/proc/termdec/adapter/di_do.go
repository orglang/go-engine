package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/descsem"
	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

var Package = do.Package(
	do.Lazy[termdec.Repo](func(i do.Injector) (termdec.Repo, error) {
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
	do.LazyNamed[descsem.Repo]("proc-descsem", func(i do.Injector) (descsem.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return descsemadapter.NewDaoPgx("proc_desc_binds ")(logger), nil
	}),
	do.Lazy[termdec.API](func(i do.Injector) (termdec.API, error) {
		termDecRepo, err := do.Invoke[termdec.Repo](i)
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
		return termdec.NewService(termDecRepo, descRootRepo, descSemRepo, typeDefRepo, transactor, logger), nil
	}),
	do.Lazy[*Controller](func(i do.Injector) (*Controller, error) {
		api, err := do.Invoke[termdec.API](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return NewControllerEcho(api, logger), nil
	}),
	do.Lazy[Presenter](func(i do.Injector) (Presenter, error) {
		api, err := do.Invoke[termdec.API](i)
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
