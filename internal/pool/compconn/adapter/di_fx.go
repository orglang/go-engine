package adapter

import (
	"go.uber.org/fx"

	compconn "orglang/go-engine/internal/pool/compconn/core"
)

var Module = fx.Module("pool/compconn",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(compconn.Repo))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(compconn.QueryBuilder))),
	),
)
