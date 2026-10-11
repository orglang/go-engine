package adapter

import (
	"go.uber.org/fx"

	"orglang/go-engine/adt/descsem"

	typedef "orglang/go-engine/pool/typedef/core"
)

var Module = fx.Module("pool/typedef",
	fx.Provide(
		fx.Annotate(typedef.NewService, fx.As(new(typedef.API))),
		fx.Annotate(NewDaoPgx, fx.As(new(typedef.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		fx.Annotate(newSQLBuilder, fx.As(new(typedef.QueryBuilder))),
		fx.Annotate(descsem.NewDaoPgx(descBinds), fx.As(new(descsem.Repo))),
	),
	fx.Invoke(
		cfgEchoController,
	),
)
