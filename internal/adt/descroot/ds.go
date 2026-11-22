package descroot

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

// Repo is the persistence port for Description aggregate revision tracking.
// A successful TouchRef advances the aggregate revision only when the
// supplied revision matches the currently persisted revision.
type Repo interface {
	TouchRef(context.Context, tx.Conn, SemRef) error
}
