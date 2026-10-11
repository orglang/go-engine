package adapter

import (
	"go.uber.org/fx"

	termdef "orglang/go-engine/proc/termdef/core"
)

var Module = fx.Module("proc/termdef",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(termdef.Repo))),
	),
)
