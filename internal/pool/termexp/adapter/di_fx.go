package adapter

import (
	"go.uber.org/fx"

	termexp "orglang/go-engine/internal/pool/termexp/core"
)

var Module = fx.Module("pool/termexp",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(termexp.Repo))),
	),
)
