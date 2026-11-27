package bootstrap

import (
	"context"
	"errors"

	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/app"
	"orglang/go-engine/internal/app/web"
	"orglang/go-engine/internal/lib/db"
	"orglang/go-engine/internal/lib/kv"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/lib/ws"
	"orglang/go-engine/internal/pool"
	poolexecevaladapter "orglang/go-engine/internal/pool/execeval/adapter"
	"orglang/go-engine/internal/proc"
	progadapter "orglang/go-engine/internal/prog/adapter"
)

// New builds the production application with samber/do. It prepares routes
// and subscriptions but leaves resource startup to the caller.
func New() (*app.App, error) {
	return NewWithPackages(
		lf.Package,
		kv.Package,
		db.Package,
		ws.Package,
		pool.Package,
		proc.Package,
		progadapter.Package,
		web.Package,
	)
}

// NewWithPackages composes an application from the supplied packages.
// It is also the composition seam used by tests and alternative entrypoints.
func NewWithPackages(packages ...func(do.Injector)) (*app.App, error) {
	application := app.NewApp(packages...)
	injector := application.Container()

	// Resolve the transactor first so database readiness is registered before
	// the HTTP listener's lifecycle hook.
	if _, err := do.Invoke[tx.Transactor](injector); err != nil {
		return stopAfterError(application, err)
	}

	server, err := do.Invoke[*echo.Echo](injector)
	if err != nil {
		return stopAfterError(application, err)
	}

	webController, err := do.Invoke[web.Controller](injector)
	if err != nil {
		return stopAfterError(application, err)
	}
	web.RegisterRoutes(server, webController)

	if err := pool.RegisterRoutes(server, injector); err != nil {
		return stopAfterError(application, err)
	}
	if err := proc.RegisterRoutes(server, injector); err != nil {
		return stopAfterError(application, err)
	}

	programController, err := do.Invoke[progadapter.Controller](injector)
	if err != nil {
		return stopAfterError(application, err)
	}
	if err := progadapter.RegisterRoutes(server, programController); err != nil {
		return stopAfterError(application, err)
	}

	if err := poolexecevaladapter.SubscribeBroker(injector); err != nil {
		return stopAfterError(application, err)
	}

	return application, nil
}

func stopAfterError(application *app.App, cause error) (*app.App, error) {
	return nil, errors.Join(cause, application.Stop(context.Background()))
}
