package adapter

import (
	"go.uber.org/fx"

	"orglang/go-engine/internal/adt/typesem"
	typesemadapter "orglang/go-engine/internal/adt/typesem/adapter"
	webapi "orglang/go-engine/internal/app/web/api"
	"orglang/go-engine/internal/lib/te"

	termdec "orglang/go-engine/internal/proc/termdec/core"
)

var Module = fx.Module("proc/termdec",
	fx.Provide(
		fx.Annotate(termdec.NewService, fx.As(new(termdec.API))),
		fx.Annotate(NewDaoPgx, fx.As(new(termdec.Repo))),
	),
	fx.Provide(
		fx.Annotate(newControllerEcho, fx.As(fx.Self(), new(webapi.ServerInterface))),
	),
	fx.Provide(
		fx.Private,
		newEchoPresenter,
		fx.Annotate(newSQLBuilder, fx.As(new(termdec.QueryBuilder))),
		fx.Annotate(newRendererStdlib, fx.As(new(te.Renderer))),
		fx.Annotate(typesemadapter.NewDaoPgx(typeDefs, descBinds), fx.As(new(typesem.Repo))),
	),
	fx.Invoke(
		cfgEchoController,
		cfgGeneratedServer,
		cfgEchoPresenter,
	),
)
