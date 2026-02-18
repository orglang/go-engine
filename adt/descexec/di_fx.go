package descexec

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/descexec",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
