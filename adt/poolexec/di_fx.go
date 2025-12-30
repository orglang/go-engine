//go:build !goverter

package poolexec

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/poolexec",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
	),
	fx.Provide(
		fx.Private,
		newHandlerEcho,
		newStepHandlerEcho,
		fx.Annotate(newDaoPgx, fx.As(new(Repo))),
	),
	fx.Invoke(
		cfgHandlerEcho,
		cfgStepHandlerEcho,
	),
)
