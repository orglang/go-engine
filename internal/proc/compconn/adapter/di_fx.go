package adapter

import (
	"go.uber.org/fx"

	compconn "orglang/go-engine/internal/proc/compconn/core"
)

var Module = fx.Module("proc/compconn",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(compconn.Repo))),
	),
)
