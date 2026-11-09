package adapter

import (
	"go.uber.org/fx"

	compconncore "orglang/go-engine/internal/proc/compconn/core"
)

var Module = fx.Module("proc/compconn",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(compconncore.Repo))),
	),
)
