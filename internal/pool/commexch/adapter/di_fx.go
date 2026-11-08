package adapter

import (
	"go.uber.org/fx"

	commexch "orglang/go-engine/pool/commexch/core"
)

var Module = fx.Module("pool/commexch",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(commexch.Repo))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(commexch.QueryBuilder))),
	),
)
