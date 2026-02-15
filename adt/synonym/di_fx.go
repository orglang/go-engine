package synonym

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/synonym",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
