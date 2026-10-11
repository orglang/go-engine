package adapter

import (
	"go.uber.org/fx"

	"orglang/go-engine/lib/te"

	"orglang/go-engine/adt/descsem"

	proctypedef "orglang/go-engine/proc/typedef/core"
)

var Module = fx.Module("adproct/typedef",
	fx.Provide(
		fx.Annotate(proctypedef.NewService, fx.As(new(proctypedef.API))),
		fx.Annotate(NewDaoPgx, fx.As(new(proctypedef.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		newEchoPresenter,
		fx.Annotate(newSQLBuilder, fx.As(new(proctypedef.QueryBuilder))),
		fx.Annotate(newRendererStdlib, fx.As(new(te.Renderer))),
		fx.Annotate(descsem.NewDaoPgx(descBinds), fx.As(new(descsem.Repo))),
	),
	fx.Invoke(
		cfgEchoController,
		cfgEchoPresenter,
	),
)
