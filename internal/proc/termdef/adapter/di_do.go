package adapter

import (
	"log/slog"
	"github.com/samber/do/v2"
	webapi "orglang/go-engine/internal/app/web/api"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	termdecadapter "orglang/go-engine/internal/proc/termdec/adapter"
	"orglang/go-engine/internal/proc/termdef/core"
)

var Package = do.Package(
	do.Lazy[core.Repo](func(i do.Injector) (core.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil { return nil, err }
		client, err := do.Invoke[*ent.Client](i)
		if err != nil { return nil, err }
		return NewDaoEnt(client, logger), nil
	}),
	do.Lazy[core.API](func(i do.Injector) (core.API, error) {
		termDecRepo, err := do.Invoke[termdec.Repo](i)
		if err != nil { return nil, err }
		termDefRepo, err := do.Invoke[core.Repo](i)
		if err != nil { return nil, err }
		transactor, err := do.Invoke[tx.Transactor](i)
		if err != nil { return nil, err }
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil { return nil, err }
		return core.NewService(termDecRepo, termDefRepo, transactor, logger), nil
	}),
	do.Lazy[webapi.ServerInterface](func(i do.Injector) (webapi.ServerInterface, error) {
		decController, err := do.Invoke[*termdecadapter.Controller](i)
		if err != nil { return nil, err }
		defAPI, err := do.Invoke[core.API](i)
		if err != nil { return nil, err }
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil { return nil, err }
		return newControllerEcho(decController, defAPI, logger), nil
	}),
)
