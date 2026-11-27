package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/uniqsym"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, EvalRec) error
	ModifyRec(context.Context, tx.Conn, EvalMod) error
	GetSnapByRef(context.Context, tx.Conn, comproot.RootRef) (EvalSnap2, error)
	GetSnapMapByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]EvalSnap1, error)
}
