package adapter

import (
	"go.uber.org/fx"

	"orglang/go-engine/adt/descsem"

	termdef "orglang/go-engine/pool/termdef/core"
)

var Module = fx.Module("pool/termdef",
	fx.Provide(
		fx.Annotate(termdef.NewService, fx.As(new(termdef.API))),
		fx.Annotate(newDaoPgx, fx.As(new(termdef.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		fx.Annotate(newSQLBuilder, fx.As(new(termdef.QueryBuilder))),
		fx.Annotate(descsem.NewDaoPgx(descBinds), fx.As(new(descsem.Repo))),
	),
	fx.Invoke(
		cfgEchoController,
	),
)
