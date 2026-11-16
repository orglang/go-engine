package pool

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descsem"
	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"
	"orglang/go-engine/internal/lib/wp"
	compconn "orglang/go-engine/internal/pool/compconn/adapter"
	compexec "orglang/go-engine/internal/pool/compexec/adapter"
	compvar "orglang/go-engine/internal/pool/compvar/adapter"
	termdef "orglang/go-engine/internal/pool/termdef/adapter"
	typedef "orglang/go-engine/internal/pool/typedef/adapter"
	typeexp "orglang/go-engine/internal/pool/typeexp/adapter"
)

// Package is the pool-domain composition boundary. It registers persistence
// adapters, services, and HTTP controllers; route registration stays explicit.
// Shared process providers are registered by proc.Package.
var Package = do.Package(
	wp.Package,
	compconn.Package,
	compvar.Package,
	compexec.Package,
	typeexp.Package,
	typedef.Package,
	termdef.Package,
	do.Lazy[descsem.Repo](func(i do.Injector) (descsem.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return descsemadapter.NewDaoPgx("pool_desc_binds ")(logger), nil
	}),
)

// RegisterRoutes mounts the pool HTTP adapters on the application's Echo server.
func RegisterRoutes(server *echo.Echo, i do.Injector) error {
	execController, err := do.Invoke[compexec.Controller](i)
	if err != nil {
		return err
	}
	if err := compexec.RegisterRoutes(server, execController); err != nil {
		return err
	}

	typeController, err := do.Invoke[typedef.Controller](i)
	if err != nil {
		return err
	}
	if err := typedef.RegisterRoutes(server, typeController); err != nil {
		return err
	}

	termController, err := do.Invoke[termdef.Controller](i)
	if err != nil {
		return err
	}
	return termdef.RegisterRoutes(server, termController)
}
