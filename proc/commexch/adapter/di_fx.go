package adapter

import (
	"go.uber.org/fx"

	commexch "orglang/go-engine/proc/commexch/core"
)

var Module = fx.Module("proc/commexch",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(commexch.Repo))),
	),
)
