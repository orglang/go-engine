package compsem

import (
	"context"
	"orglang/go-engine/lib/tx"
)

type Repo interface {
	TouchRef(context.Context, tx.Conn, SemRef) error
}

type SemRefDS struct {
	CompID string `db:"comp_id"`
	CompRN int64  `db:"comp_rn"`
}
