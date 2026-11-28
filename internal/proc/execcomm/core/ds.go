package core

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	ModifyRec(context.Context, tx.Conn, CommMod) error
	GetSnapByQry(context.Context, tx.Conn, CommQry) (CommSnap, error)
}
