package execroot

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

// Repo is the persistence port for execution aggregate revision tracking.
type Repo interface {
	TouchRef(context.Context, tx.Conn, RootRef) error
}
