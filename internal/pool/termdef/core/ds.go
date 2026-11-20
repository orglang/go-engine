package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/uniqsym"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DefRec) error
	GetRecByQN(context.Context, tx.Conn, uniqsym.ADT) (DefRec, error)
}
