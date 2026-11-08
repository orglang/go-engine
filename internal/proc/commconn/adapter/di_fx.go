package adapter

import (
	"go.uber.org/fx"

	commconn "orglang/go-engine/internal/proc/commconn/core"
)

var Module = fx.Module("proc/commconn",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(commconn.Repo))),
	),
)
