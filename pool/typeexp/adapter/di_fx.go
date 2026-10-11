package adapter

import (
	"go.uber.org/fx"

	typeexp "orglang/go-engine/pool/typeexp/core"
)

var Module = fx.Module("pool/typeexp",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(typeexp.Repo))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(typeexp.QueryBuilder))),
	),
)
