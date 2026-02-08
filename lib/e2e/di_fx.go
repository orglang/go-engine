package e2e

import (
	"go.uber.org/fx"
)

var Module = fx.Module("lib/e2e",
	fx.Provide(
		newPoolExecAPI,
		newPoolDecAPI,
		newXactDefAPI,
		newProcDecAPI,
		newProcExecAPI,
		newTypeDefAPI,
	),
)
