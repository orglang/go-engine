package adapter

import (
	"go.uber.org/fx"

	"orglang/go-engine/adt/descsem"

	pooltypedef "orglang/go-engine/pool/typedef/core"
)

var Module = fx.Module("pool/typedef",
	fx.Provide(
		fx.Annotate(pooltypedef.NewService, fx.As(new(pooltypedef.API))),
		fx.Annotate(NewDaoPgx, fx.As(new(pooltypedef.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		fx.Annotate(newSQLBuilder, fx.As(new(pooltypedef.QueryBuilder))),
		fx.Annotate(descsem.NewDaoPgx(descBinds), fx.As(new(descsem.Repo))),
	),
	fx.Invoke(
		cfgEchoController,
	),
)
