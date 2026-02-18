package descbind

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/descbind",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
