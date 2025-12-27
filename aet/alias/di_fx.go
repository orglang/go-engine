//go:build !goverter

package alias

import (
	"go.uber.org/fx"
)

var Module = fx.Module("aet/alias",
	fx.Provide(
		fx.Annotate(newRepoPgx, fx.As(new(Repo))),
	),
)
