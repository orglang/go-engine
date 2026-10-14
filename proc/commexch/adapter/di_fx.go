package adapter

import (
	"go.uber.org/fx"

	corecommexch "orglang/go-engine/proc/commexch/core"
)

var Module = fx.Module("proc/commexch",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(corecommexch.Repo))),
	),
)
