//go:build !goverter

package procxec

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/procxec",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
	),
	fx.Provide(
		fx.Private,
		newHandlerEcho,
		fx.Annotate(newDaoPgx, fx.As(new(execRepo))),
	),
	fx.Invoke(
		cfgHandlerEcho,
	),
)
