package adapter

import (
	"go.uber.org/fx"

	commturn "orglang/go-engine/pool/commturn/core"
)

var Module = fx.Module("pool/commturn",
	fx.Provide(
		fx.Annotate(commturn.NewService, fx.As(new(commturn.API))),
		fx.Annotate(NewDaoPgx, fx.As(new(commturn.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		fx.Annotate(newSQLBuilder, fx.As(new(commturn.QueryBuilder))),
	),
	fx.Invoke(
		cfgEchoController,
	),
)
