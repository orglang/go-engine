package adapter

import (
	"go.uber.org/fx"

	pooltermexp "orglang/go-engine/pool/termexp/core"
)

var Module = fx.Module("pool/termexp",
	fx.Provide(
		fx.Annotate(NewDaoPgx, fx.As(new(pooltermexp.Repo))),
	),
)
