package poolstep

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/poolstep",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
	),
	fx.Provide(
		fx.Private,
		newEchoController,
	),
	fx.Invoke(
		cfgEchoController,
	),
)
