package wp

import (
	"go.uber.org/fx"
)

var Module = fx.Module("lib/wp",
	fx.Provide(
		newPondPool,
		newWorkerPool,
	),
)
