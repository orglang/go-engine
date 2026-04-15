package termexp

import (
	"go.uber.org/fx"
)

var Module = fx.Module("proc/termexp",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
