package procconn

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/procconn",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
