//go:build !goverter

package typeexp

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/typeexp",
	fx.Provide(
		fx.Annotate(newDaoPgx, fx.As(new(Repo))),
	),
)
