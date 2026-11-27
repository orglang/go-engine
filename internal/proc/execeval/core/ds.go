package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/comproot"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, EvalRec) error
	ModifyRec(context.Context, tx.Conn, EvalMod) error
	GetSnapByRef(context.Context, tx.Conn, comproot.RootRef) (EvalSnap, error)
}
