package adapter

import (
	"go.uber.org/fx"

	commconn "orglang/go-engine/internal/pool/commconn/core"
)

var Module = fx.Module("pool/commconn",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(commconn.Repo))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(commconn.QueryBuilder))),
	),
)
