package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/execroot"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, EvalRec) error
	ModifyRec(context.Context, tx.Conn, EvalMod) error
	GetSnapByRef(context.Context, tx.Conn, execroot.RootRef) (EvalSnap, error)
}
