package adapter

import (
	"go.uber.org/fx"

	coreprog "orglang/go-engine/internal/prog/core"
)

var Module = fx.Module("prog/adapter",
	fx.Provide(
		fx.Annotate(coreprog.NewService, fx.As(new(coreprog.API))),
		newControllerEcho,
	),
	fx.Invoke(
		cfgGeneratedServer,
	),
)
