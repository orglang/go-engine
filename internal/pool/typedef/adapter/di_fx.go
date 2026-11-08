package adapter

import (
	"go.uber.org/fx"

	"orglang/go-engine/internal/adt/descsem"
	descsemadapter "orglang/go-engine/internal/adt/descsem/adapter"

	typedef "orglang/go-engine/internal/pool/typedef/core"
)

var Module = fx.Module("pool/typedef",
	fx.Provide(
		fx.Annotate(typedef.NewService, fx.As(new(typedef.API))),
		fx.Annotate(NewDaoPgx, fx.As(new(typedef.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		fx.Annotate(newSQLBuilder, fx.As(new(typedef.QueryBuilder))),
		fx.Annotate(descsemadapter.NewDaoPgx(descBinds), fx.As(new(descsem.Repo))),
	),
	fx.Invoke(
		cfgEchoController,
	),
)
