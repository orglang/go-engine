package poolimpl

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/poolimpl",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
		fx.Annotate(newPondBroker, fx.As(new(Exch))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(newSQLBuilder, fx.As(new(queryBuilder))),
		newEchoController,
		newPondBroker,
	),
	fx.Invoke(
		cfgEchoController,
	),
)
