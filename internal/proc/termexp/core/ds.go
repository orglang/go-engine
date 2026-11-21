package core

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	Insert(context.Context, tx.Conn, ExpRec) error
}
