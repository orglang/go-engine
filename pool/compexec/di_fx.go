package compexec

import (
	"go.uber.org/fx"
)

var Module = fx.Module("pool/compexec",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
	fx.Provide(
		fx.Private,
		newEchoController,
		newPondBroker,
		fx.Annotate(newPondBroker, fx.As(new(Exch))),
		// fx.Annotate(newWorkerPoolBroker, fx.As(new(Exch))),
		fx.Annotate(newSQLBuilder, fx.As(new(queryBuilder))),
	),
	fx.Invoke(
		cfgEchoController,
		cfgPondBroker,
		// fx.Annotate(cfgPondBroker, fx.From(new(Exch))),
		// cfgWorkerPoolBroker,
	),
)
