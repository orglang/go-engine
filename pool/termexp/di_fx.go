package termexp

import (
	"go.uber.org/fx"
)

var Module = fx.Module("pool/termexp",
	fx.Provide(
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
)
