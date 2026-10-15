package adapter

import (
	"go.uber.org/fx"

	compexec "orglang/go-engine/proc/compexec/core"
)

var Module = fx.Module("proc/compexec",
	fx.Provide(
		fx.Annotate(compexec.NewService, fx.As(new(compexec.API))),
		fx.Annotate(newDaoPgx, fx.As(new(compexec.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		fx.Annotate(newSQLBuilder, fx.As(new(compexec.QueryBuilder))),
	),
	fx.Invoke(
		cfgEchoController,
	),
)
