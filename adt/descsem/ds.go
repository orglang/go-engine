package descsem

import (
	"context"
	"orglang/go-engine/lib/tx"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, SemRec) error
}

type SemRecDS struct {
	DescQN string `db:"desc_qn"`
	DescID string `db:"desc_id"`
	Kind   int16  `db:"kind"`
}
