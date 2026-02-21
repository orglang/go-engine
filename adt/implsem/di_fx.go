package implsem

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/implsem",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
