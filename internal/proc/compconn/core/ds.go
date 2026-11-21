package core

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	ModifyRec(context.Context, tx.Conn, ConnMod) error
	GetSnapByQry(context.Context, tx.Conn, ConnQry) (ConnSnap, error)
}
