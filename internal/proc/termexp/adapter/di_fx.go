package adapter

import (
	"go.uber.org/fx"

	termexp "orglang/go-engine/proc/termexp/core"
)

var Module = fx.Module("proc/termexp",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(termexp.Repo))),
	),
)
