//go:build !goverter

package typealias

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/typealias",
	fx.Provide(
		fx.Annotate(newDaoPgx, fx.As(new(Repo))),
	),
)
