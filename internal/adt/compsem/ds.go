package compsem

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	TouchRef(context.Context, tx.Conn, SemRef) error
}
