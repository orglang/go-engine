package comproot

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

// Repo is the persistence port for Computation aggregate revision tracking.
type Repo interface {
	TouchRef(context.Context, tx.Conn, RootRef) error
}
