package descroot

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

// Repo persists Description aggregate references and enforces their OCC boundary.
type Repo interface {
	AddRef(context.Context, tx.Conn, SemRef) error
	TouchRef(context.Context, tx.Conn, SemRef) error
}
