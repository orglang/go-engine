package adapter

import (
	"go.uber.org/fx"

	"orglang/go-engine/adt/typesem"
	"orglang/go-engine/lib/te"

	termdec "orglang/go-engine/proc/termdec/core"
)

var Module = fx.Module("proc/termdec",
	fx.Provide(
		fx.Annotate(termdec.NewService, fx.As(new(termdec.API))),
		fx.Annotate(NewDaoPgx, fx.As(new(termdec.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		newEchoPresenter,
		fx.Annotate(newSQLBuilder, fx.As(new(termdec.QueryBuilder))),
		fx.Annotate(newRendererStdlib, fx.As(new(te.Renderer))),
		fx.Annotate(typesem.NewDaoPgx(typeDefs, descBinds), fx.As(new(typesem.Repo))),
	),
	fx.Invoke(
		cfgEchoController,
		cfgEchoPresenter,
	),
)
