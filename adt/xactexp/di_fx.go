package xactexp

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/xactexp",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
