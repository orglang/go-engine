package poolstep

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/poolstep",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
	fx.Provide(
		fx.Private,
		newEchoController,
	),
	fx.Invoke(
		cfgEchoController,
	),
)
