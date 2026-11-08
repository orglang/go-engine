package adapter

import (
	"go.uber.org/fx"

	"orglang/go-engine/internal/adt/compsem"
	compsemadapter "orglang/go-engine/internal/adt/compsem/adapter"
	"orglang/go-engine/internal/adt/implsem"
	implsemadapter "orglang/go-engine/internal/adt/implsem/adapter"

	compexec "orglang/go-engine/internal/pool/compexec/core"
)

var Module = fx.Module("pool/compexec",
	fx.Provide(
		fx.Annotate(compexec.NewService, fx.As(new(compexec.API))),
		fx.Annotate(NewDaoPgx, fx.As(new(compexec.Repo))),
	),
	fx.Provide(
		fx.Private,
		newControllerEcho,
		newPondBroker,
		fx.Annotate(newPondBroker, fx.As(new(compexec.Broker))),
		// fx.Annotate(newWorkerPoolBroker, fx.As(new(compexec.Broker))),
		fx.Annotate(newSQLBuilder, fx.As(new(compexec.QueryBuilder))),
		fx.Annotate(implsemadapter.NewDaoPgx(implBinds), fx.As(new(implsem.Repo))),
		fx.Annotate(compsemadapter.NewDaoPgx(compExecs), fx.As(new(compsem.Repo))),
	),
	fx.Invoke(
		cfgEchoController,
		cfgPondBroker,
		// fx.Annotate(cfgPondBroker, fx.From(new(compexec.Broker))),
		// cfgWorkerPoolBroker,
	),
)
