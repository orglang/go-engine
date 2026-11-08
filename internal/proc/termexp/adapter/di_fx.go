package adapter

import (
	"go.uber.org/fx"

	termexp "orglang/go-engine/internal/proc/termexp/core"
)

var Module = fx.Module("proc/termexp",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(termexp.Repo))),
	),
)
