package poolvar

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/poolvar",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(queryBuilder))),
	),
)
