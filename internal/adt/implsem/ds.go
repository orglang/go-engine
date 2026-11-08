package implsem

import (
	"context"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, SemRec) error
}

type SemRecDS struct {
	ImplQN string `db:"impl_qn"`
	ImplID string `db:"impl_id"`
	Kind   int16  `db:"kind"`
}
