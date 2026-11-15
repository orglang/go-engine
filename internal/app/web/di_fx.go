package web

import (
	"go.uber.org/fx"

	"orglang/go-engine/internal/lib/te"
	webapi "orglang/go-engine/internal/app/web/api"
)

var Module = fx.Module("app/web",
	fx.Provide(
		fx.Private,
		fx.Annotate(newRendererStdlib, fx.As(new(te.Renderer))),
		newControllerEcho,
	),
	fx.Invoke(
		cfgEchoController,
		cfgGeneratedServer,
	),
)
