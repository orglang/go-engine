package poolexec

import (
	"go.uber.org/fx"

	"orglang/go-engine/lib/te"
)

var Module = fx.Module("adt/poolexec",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
	fx.Provide(
		fx.Private,
		newEchoController,
		fx.Annotate(newRendererStdlib, fx.As(new(te.Renderer))),
	),
	fx.Invoke(
		cfgEchoController,
	),
)
