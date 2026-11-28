package core

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, CommRec) error
	ModifyRec(context.Context, tx.Conn, CommMod) error
	GetSnapByQry(context.Context, tx.Conn, CommQry) (CommSnap, error)
}
