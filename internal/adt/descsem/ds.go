package descsem

import (
	"context"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, SemRec) error
}
