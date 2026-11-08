package adapter

import (
	"go.uber.org/fx"

	"orglang/go-engine/internal/adt/descsem"
	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"

	termdef "orglang/go-engine/internal/pool/termdef/core"
)

var Module = fx.Module("pool/termdef",
	fx.Provide(
		fx.Annotate(termdef.NewService, fx.As(new(termdef.API))),
		fx.Annotate(NewDaoPgx, fx.As(new(termdef.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		fx.Annotate(newSQLBuilder, fx.As(new(termdef.QueryBuilder))),
		fx.Annotate(descsemadapter.NewDaoPgx(descBinds), fx.As(new(descsem.Repo))),
	),
	fx.Invoke(
		cfgEchoController,
	),
)
