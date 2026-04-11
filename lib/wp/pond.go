package wp

import (
	"github.com/alitto/pond/v2"
)

func newPondPool() pond.Pool {
	return pond.NewPool(10, pond.WithoutPanicRecovery())
}
