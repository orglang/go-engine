package proc

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descsem"
	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"
	"orglang/go-engine/internal/adt/typesem"
	typesemadapter "orglang/go-engine/internal/adt/typesem/adapter"
	webapi "orglang/go-engine/internal/app/web/api"
	compconn "orglang/go-engine/internal/proc/compconn/adapter"
	compexec "orglang/go-engine/internal/proc/compexec/adapter"
	termdec "orglang/go-engine/internal/proc/termdec/adapter"
	termdef "orglang/go-engine/internal/proc/termdef/adapter"
	typedef "orglang/go-engine/internal/proc/typedef/adapter"
	typeexp "orglang/go-engine/internal/proc/typeexp/adapter"
)

// Package is the process-domain composition boundary. It registers the
// persistence adapters and services currently used by production process
// endpoints. Incomplete and unused areas such as termexp remain unwired.
var Package = do.Package(
	compconn.Package,
	compexec.Package,
	termdec.Package,
	termdef.Package,
	typedef.Package,
	typeexp.Package,
	do.Lazy[typesem.Repo](func(i do.Injector) (typesem.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return typesemadapter.NewDaoPgx("proc_type_defs ", "proc_desc_binds ")(logger), nil
	}),
	do.Lazy[descsem.Repo](func(i do.Injector) (descsem.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return descsemadapter.NewDaoPgx("proc_desc_binds ")(logger), nil
	}),
)

// RegisterRoutes resolves the process HTTP adapters and mounts the currently
// migrated API and server-rendered routes on the application's Echo server.
// It is an explicit composition function; it does not start an HTTP listener.
func RegisterRoutes(server *echo.Echo, i do.Injector) error {
	execController, err := do.Invoke[compexec.Controller](i)
	if err != nil {
		return err
	}
	termDecController, err := do.Invoke[webapi.ServerInterface](i)
	if err != nil {
		return err
	}
	termDecPresenter, err := do.Invoke[termdec.Presenter](i)
	if err != nil {
		return err
	}
	typeDefController, err := do.Invoke[typedef.Controller](i)
	if err != nil {
		return err
	}
	typeDefPresenter, err := do.Invoke[typedef.Presenter](i)
	if err != nil {
		return err
	}

	if err := compexec.RegisterRoutes(server, execController); err != nil {
		return err
	}
	if err := termdec.RegisterGeneratedRoutes(server, termDecController); err != nil {
		return err
	}
	if err := termdec.RegisterPresenterRoutes(server, termDecPresenter); err != nil {
		return err
	}
	if err := typedef.RegisterRoutes(server, typeDefController); err != nil {
		return err
	}
	return typedef.RegisterPresenterRoutes(server, typeDefPresenter)
}
