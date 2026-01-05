//go:build !goverter

package procexp

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/procexp",
	fx.Provide(
		fx.Annotate(newDaoPgx, fx.As(new(Repo))),
	),
)
