package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/compsem"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ExecRec) error
	ModifyRec(context.Context, tx.Conn, ExecMod) error
	GetSnapByRef(context.Context, tx.Conn, compsem.SemRef) (ExecSnap, error)
}
