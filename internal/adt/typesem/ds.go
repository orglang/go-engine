package typesem

import (
	"context"

	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	TouchRef(context.Context, tx.Conn, SemRef) error
	GetRefsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]SemRef, error)
}
