package adapter

import (
	"go.uber.org/fx"

	proccommturn "orglang/go-engine/proc/commturn/core"
)

var Module = fx.Module("proc/commturn",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(proccommturn.Repo))),
	),
)
