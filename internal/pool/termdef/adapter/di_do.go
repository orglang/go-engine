package adapter

import (
	"log/slog"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	"orglang/go-engine/internal/pool/termdef/core"
	termexp "orglang/go-engine/internal/pool/termexp/core"
	typedef "orglang/go-engine/internal/pool/typedef/core"
)

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
		return NewDaoEnt(client, logger), nil
	}),
	do.Lazy[core.API](func(i do.Injector) (core.API, error) {
		repo, err := do.Invoke[core.Repo](i)
		if err != nil {
			return nil, err
		}
		termExpRepo, err := do.Invoke[termexp.Repo](i)
		if err != nil {
			return nil, err
		}
		descRootRepo, err := do.InvokeNamed[descroot.Repo](i, "pool-descroot")
		if err != nil {
			return nil, err
		}
		typeDefRepo, err := do.Invoke[typedef.Repo](i)
		if err != nil {
			return nil, err
		}
		descRepo, err := do.InvokeNamed[descsem.Repo](i, "pool-descsem")
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
		return core.NewService(repo, termExpRepo, descRootRepo, typeDefRepo, descRepo, transactor, logger), nil
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
