package adapter

import (
	"go.uber.org/fx"

	conncore "orglang/go-engine/internal/proc/conn/core"
)

var Module = fx.Module("proc/conn",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(conncore.Repo))),
	),
)
