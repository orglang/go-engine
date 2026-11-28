package pool

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/progroot"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/wp"
	"orglang/go-engine/internal/pool/ent"
	execcomm "orglang/go-engine/internal/pool/execcomm/adapter"
	execeval "orglang/go-engine/internal/pool/execeval/adapter"
	execvar "orglang/go-engine/internal/pool/execvar/adapter"
	progrootadapter "orglang/go-engine/internal/pool/progroot/adapter"
	termdec "orglang/go-engine/internal/pool/termdec/adapter"
	termdef "orglang/go-engine/internal/pool/termdef/adapter"
	typedef "orglang/go-engine/internal/pool/typedef/adapter"
)

var Package = do.Package(
	do.Lazy[*ent.Client](func(do.Injector) (*ent.Client, error) {
		return ent.NewClient(ent.Driver(entpgx.New())), nil
	}),
	wp.Package,
	execcomm.Package,
	execvar.Package,
	execeval.Package,
	typedef.Package,
	termdec.Package,
	termdef.Package,
	do.LazyNamed[progroot.Repo]("pool-progroot", func(i do.Injector) (progroot.Repo, error) {
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
)

func RegisterRoutes(server *echo.Echo, i do.Injector) error {
	evalController, err := do.Invoke[execeval.Controller](i)
	if err != nil {
		return err
	}
	if err := execeval.RegisterRoutes(server, evalController); err != nil {
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
