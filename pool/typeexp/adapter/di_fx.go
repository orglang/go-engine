package adapter

import (
	"go.uber.org/fx"

	pooltypeexp "orglang/go-engine/pool/typeexp/core"
)

var Module = fx.Module("pool/typeexp",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(pooltypeexp.Repo))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(pooltypeexp.QueryBuilder))),
	),
)
