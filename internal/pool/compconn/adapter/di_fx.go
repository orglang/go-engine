package adapter

import (
	"go.uber.org/fx"

	compconncore "orglang/go-engine/internal/pool/compconn/core"
)

var Module = fx.Module("pool/compconn",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(compconncore.Repo))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(compconncore.QueryBuilder))),
	),
)
