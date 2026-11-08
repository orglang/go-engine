package web

import (
	"go.uber.org/fx"

	"orglang/go-engine/internal/lib/te"
)

var Module = fx.Module("app/web",
	fx.Provide(
		fx.Private,
		fx.Annotate(newRendererStdlib, fx.As(new(te.Renderer))),
		newControllerEcho,
	),
	fx.Invoke(
		cfgEchoController,
	),
)
