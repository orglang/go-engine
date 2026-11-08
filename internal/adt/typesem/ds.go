package typesem

import (
	"context"
	"orglang/go-engine/lib/tx"

	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	TouchRef(context.Context, tx.Conn, SemRef) error
	GetRefsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]SemRef, error)
}

type SemRefDS struct {
	TypeID string `db:"type_id"`
	TypeRN int64  `db:"type_rn"`
}
