package termdef

import (
	"go.uber.org/fx"
)

var Module = fx.Module("proc/termdef",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
