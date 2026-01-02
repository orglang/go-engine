//go:build !goverter

package termsyn

import (
	"go.uber.org/fx"
)

var Module = fx.Module("adt/termsyn",
	fx.Provide(
		fx.Annotate(newDaoPgx, fx.As(new(Repo))),
	),
)
