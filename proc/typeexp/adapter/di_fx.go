package adapter

import (
	"go.uber.org/fx"

	proctypeexp "orglang/go-engine/proc/typeexp/core"
)

var Module = fx.Module("proc/typeexp",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(proctypeexp.Repo))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(proctypeexp.QueryBuilder))),
	),
)
