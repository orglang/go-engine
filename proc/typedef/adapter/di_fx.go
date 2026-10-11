package adapter

import (
	"go.uber.org/fx"

	"orglang/go-engine/lib/te"

	"orglang/go-engine/adt/descsem"

	typedef "orglang/go-engine/proc/typedef/core"
)

var Module = fx.Module("adproct/typedef",
	fx.Provide(
		fx.Annotate(typedef.NewService, fx.As(new(typedef.API))),
		fx.Annotate(NewDaoPgx, fx.As(new(typedef.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		newEchoPresenter,
		fx.Annotate(newSQLBuilder, fx.As(new(typedef.QueryBuilder))),
		fx.Annotate(newRendererStdlib, fx.As(new(te.Renderer))),
		fx.Annotate(descsem.NewDaoPgx(descBinds), fx.As(new(descsem.Repo))),
	),
	fx.Invoke(
		cfgEchoController,
		cfgEchoPresenter,
	),
)
