package core

import (
	"context"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/uniqsym"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DefRec) error
	GetRecsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

type DefRecDS struct {
	TypeID string `db:"type_id"`
	TypeRN int64  `db:"type_rn"`
	ExpCK  int64  `db:"exp_ck"`
}
