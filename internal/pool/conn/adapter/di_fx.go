package adapter

import (
	"go.uber.org/fx"

	conncore "orglang/go-engine/internal/pool/conn/core"
)

var Module = fx.Module("pool/conn",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(conncore.Repo))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(conncore.QueryBuilder))),
	),
)
