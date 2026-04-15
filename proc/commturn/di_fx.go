package commturn

import (
	"go.uber.org/fx"
)

var Module = fx.Module("proc/commturn",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
