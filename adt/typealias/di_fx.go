//go:build !goverter

package typealias

import (
	"go.uber.org/fx"
)

var Module = fx.Module("aet/alias",
	fx.Provide(
		fx.Annotate(newDaoPgx, fx.As(new(Repo))),
	),
)
