package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/comproot"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ExecRec) error
	ModifyRec(context.Context, tx.Conn, ExecMod) error
	GetSnapByRef(context.Context, tx.Conn, comproot.SemRef) (ExecSnap, error)
}
