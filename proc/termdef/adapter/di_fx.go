package adapter

import (
	"go.uber.org/fx"

	proctermdef "orglang/go-engine/proc/termdef/core"
)

var Module = fx.Module("proc/termdef",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(proctermdef.Repo))),
	),
)
