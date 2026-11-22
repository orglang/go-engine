package descroot

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

// Repo persists Description aggregate references and enforces their OCC boundary.
type Repo interface {
	// AddRef inserts a new aggregate root at its initial revision.
	AddRef(context.Context, tx.Conn, SemRef) error
	// TouchRef advances the persisted revision when the supplied revision matches.
	TouchRef(context.Context, tx.Conn, SemRef) error
}
