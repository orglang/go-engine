package pool

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/descsem"
	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/wp"
	compconn "orglang/go-engine/internal/pool/compconn/adapter"
	compexec "orglang/go-engine/internal/pool/compexec/adapter"
	compvar "orglang/go-engine/internal/pool/compvar/adapter"
	descrootadapter "orglang/go-engine/internal/pool/descroot/adapter"
	"orglang/go-engine/internal/pool/ent"
	termdec "orglang/go-engine/internal/pool/termdec/adapter"
	termdef "orglang/go-engine/internal/pool/termdef/adapter"
	typedef "orglang/go-engine/internal/pool/typedef/adapter"
)

var Package = do.Package(
	do.Lazy[*ent.Client](func(do.Injector) (*ent.Client, error) {
		return ent.NewClient(ent.Driver(entpgx.New())), nil
	}),
	wp.Package,
	compconn.Package,
	compvar.Package,
	compexec.Package,
	typedef.Package,
	termdec.Package,
	termdef.Package,
	do.LazyNamed[descroot.Repo]("pool-descroot", func(i do.Injector) (descroot.Repo, error) {
		client, err := do.Invoke[*ent.Client](i)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return descrootadapter.NewDaoEnt(client, logger), nil
	}),
	do.LazyNamed[descsem.Repo]("pool-descsem", func(i do.Injector) (descsem.Repo, error) {
		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		return descsemadapter.NewDaoPgx("pool_desc_binds ")(logger), nil
	}),
)

func RegisterRoutes(server *echo.Echo, i do.Injector) error {
	execController, err := do.Invoke[compexec.Controller](i)
	if err != nil {
		return err
	}
	if err := compexec.RegisterRoutes(server, execController); err != nil {
		return err
	}

	typeController, err := do.Invoke[typedef.GeneratedController](i)
	if err != nil {
		return err
	}
	if err := typedef.RegisterGeneratedRoutes(server, typeController); err != nil {
		return err
	}

	termDecController, err := do.Invoke[termdec.GeneratedController](i)
	if err != nil {
		return err
	}
	if err := termdec.RegisterGeneratedRoutes(server, termDecController); err != nil {
		return err
	}

	termDefController, err := do.Invoke[termdef.GeneratedController](i)
	if err != nil {
		return err
	}
	return termdef.RegisterGeneratedRoutes(server, termDefController)
}
