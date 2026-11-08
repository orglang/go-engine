package core

import (
	"context"
	"orglang/go-engine/lib/tx"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ExchRec) error
	Modifyec(context.Context, tx.Conn, ExchMod) error
	GetRefsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error)
	GetSnapByQry(context.Context, tx.Conn, ExchQry) (ExchSnap, error)
}
