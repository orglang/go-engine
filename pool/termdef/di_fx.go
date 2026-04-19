package termdef

import (
	"go.uber.org/fx"
)

var Module = fx.Module("pool/termdec",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(API))),
		fx.Annotate(newPgxDAO, fx.As(new(Repo))),
	),
	fx.Provide(
		fx.Private,
		newEchoController,
		fx.Annotate(newSQLBuilder, fx.As(new(queryBuilder))),
	),
	fx.Invoke(
		cfgEchoController,
	),
)
