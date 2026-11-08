package adapter

import (
	"log/slog"

	"go.uber.org/fx"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/implsem"
	implsemadapter "orglang/go-engine/internal/adt/implsem/adapter"
	compsemadapter "orglang/go-engine/internal/adt/compsem/adapter"
	commexch "orglang/go-engine/internal/pool/commexch/core"
	commturn "orglang/go-engine/internal/pool/commturn/core"
	compvarcore "orglang/go-engine/internal/pool/compvar/core"
	termdef "orglang/go-engine/internal/pool/termdef/core"
	typeexp "orglang/go-engine/internal/pool/typeexp/core"
	proccompexec "orglang/go-engine/internal/proc/compexec/core"
	"orglang/go-engine/internal/lib/tx"

	compexec "orglang/go-engine/internal/pool/compexec/core"
)


type serviceParams struct {
	fx.In

	CompExecRepo   compexec.Repo
	CompExecBroker compexec.Broker
	CompVarRepo    compvarcore.Repo
	CommExchRepo   commexch.Repo
	CommTurnRepo   commturn.Repo
	TypeExpRepo    typeexp.Repo
	ProcExecRepo   proccompexec.Repo
	TermDefRepo    termdef.Repo
	ImplSemRepo    implsem.Repo
	CompSemRepo    compsem.Repo
	Transactor     tx.Transactor
	Logger         *slog.Logger
}

func newService(params serviceParams) compexec.API {
	return compexec.NewService(compexec.ServiceDeps{
		CompExecRepo: params.CompExecRepo,
		CompExecBroker: params.CompExecBroker,
		CompVarRepo: params.CompVarRepo,
		CommExchRepo: params.CommExchRepo,
		CommTurnRepo: params.CommTurnRepo,
		TypeExpRepo: params.TypeExpRepo,
		ProcExecRepo: params.ProcExecRepo,
		TermDefRepo: params.TermDefRepo,
		ImplSemRepo: params.ImplSemRepo,
		CompSemRepo: params.CompSemRepo,
		Transactor: params.Transactor,
		Logger: params.Logger,
	})
}

var Module = fx.Module("pool/compexec",
	fx.Provide(
		fx.Annotate(newService, fx.As(new(compexec.API))),
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
