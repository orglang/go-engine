package descsem

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/descsem",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
