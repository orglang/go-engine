package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DefRec) error
	ModifyRec(context.Context, tx.Conn, DefRec) error
	GetRefs(context.Context, tx.Conn) ([]typesem.SemRef, error)
	GetRecByRef(context.Context, tx.Conn, typesem.SemRef) (DefRec, error)
	GetRecsByRefs(context.Context, tx.Conn, []typesem.SemRef) ([]DefRec, error)
	GetRecByQN(context.Context, tx.Conn, uniqsym.ADT) (DefRec, error)
	GetRecsByQNs(context.Context, tx.Conn, []uniqsym.ADT) ([]DefRec, error)
	SelectEnv(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

type DefRecDS struct {
	TypeID string `db:"type_id"`
	TypeRN int64  `db:"type_rn"`
	ExpCK  int64  `db:"exp_ck"`
}
