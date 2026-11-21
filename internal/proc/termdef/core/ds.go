package core

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	InsertProc(context.Context, tx.Conn, DefRec) error
}
