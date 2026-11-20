package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/uniqsym"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ExecRec) error
	ModifyRec(context.Context, tx.Conn, ExecMod) error
	GetSnapByRef(context.Context, tx.Conn, compsem.SemRef) (ExecSnap2, error)
	GetSnapMapByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]ExecSnap1, error)
}
