//go:build !goverter

package procdef

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/procdef",
	fx.Provide(
		fx.Annotate(newDaoPgx, fx.As(new(Repo))),
	),
)
