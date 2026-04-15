package commexch

import (
	"go.uber.org/fx"
)

var Module = fx.Module("proc/commexch",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
