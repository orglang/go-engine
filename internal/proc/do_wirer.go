package proc

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/progroot"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/proc/ent"
	execcomm "orglang/go-engine/internal/proc/execcomm/adapter"
	execeval "orglang/go-engine/internal/proc/execeval/adapter"
	progrootadapter "orglang/go-engine/internal/proc/progroot/adapter"
	termdec "orglang/go-engine/internal/proc/termdec/adapter"
	termdef "orglang/go-engine/internal/proc/termdef/adapter"
	typedef "orglang/go-engine/internal/proc/typedef/adapter"
)

// Package is the process-domain composition boundary. It registers the
// persistence adapters and services currently used by production endpoints.
var Package = do.Package(
	do.Lazy[*ent.Client](func(do.Injector) (*ent.Client, error) {
		return ent.NewClient(ent.Driver(entpgx.New())), nil
	}),
	execcomm.Package,
	do.LazyNamed[progroot.Repo]("proc-progroot", func(i do.Injector) (progroot.Repo, error) {
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return progrootadapter.NewDaoEnt(client, logger), nil
	}),
	execeval.Package,
	termdec.Package,
	termdef.Package,
	typedef.Package,
)

// RegisterRoutes resolves the process HTTP adapters and mounts the currently
// migrated API and server-rendered routes on the application's Echo server.
// It is an explicit composition function; it does not start an HTTP listener.
func RegisterRoutes(server *echo.Echo, i do.Injector) error {
	evalController, err := do.Invoke[execeval.GeneratedController](i)
	if err != nil {
		return err
	}
	termDefController, err := do.Invoke[*termdef.DefinitionController](i)
	if err != nil {
		return err
	}
	termDecPresenter, err := do.Invoke[termdec.Presenter](i)
	if err != nil {
		return err
	}
	termDecController, err := do.Invoke[*termdec.Controller](i)
	if err != nil {
		return err
	}
	typeDefController, err := do.Invoke[typedef.GeneratedController](i)
	if err != nil {
		return err
	}
	typeDefPresenter, err := do.Invoke[typedef.Presenter](i)
	if err != nil {
		return err
	}

	if err := execeval.RegisterGeneratedRoutes(server, evalController); err != nil {
		return err
	}
	if err := termdec.RegisterGeneratedRoutes(server, termDecController); err != nil {
		return err
	}
	if err := termdef.RegisterGeneratedRoutes(server, termDefController); err != nil {
		return err
	}
	if err := termdec.RegisterPresenterRoutes(server, termDecPresenter); err != nil {
		return err
	}
	if err := typedef.RegisterGeneratedRoutes(server, typeDefController); err != nil {
		return err
	}
	return typedef.RegisterPresenterRoutes(server, typeDefPresenter)
}
