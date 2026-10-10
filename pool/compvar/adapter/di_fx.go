package adapter

import (
	"go.uber.org/fx"

	poolcompvar "orglang/go-engine/pool/compvar/core"
)

var Module = fx.Module("pool/compvar",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(poolcompvar.Repo))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(poolcompvar.QueryBuilder))),
	),
)
