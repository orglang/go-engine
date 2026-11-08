package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DefRec) error
	Update(context.Context, tx.Conn, DefRec) error
	SelectRefs(context.Context, tx.Conn) ([]typesem.SemRef, error)
	SelectRecByRef(context.Context, tx.Conn, typesem.SemRef) (DefRec, error)
	SelectRecsByRefs(context.Context, tx.Conn, []typesem.SemRef) ([]DefRec, error)
	SelectRecByQN(context.Context, tx.Conn, uniqsym.ADT) (DefRec, error)
	GetRecsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

type DefRecDS struct {
	TypeID string `db:"type_id"`
	TypeRN int64  `db:"type_rn"`
	ExpVK  int64  `db:"exp_vk"`
}
