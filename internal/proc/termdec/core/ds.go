package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/termsem"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DecRec) error
	GetRefs(context.Context, tx.Conn) ([]termsem.SemRef, error)
	GetSnap(context.Context, tx.Conn, termsem.SemRef) (DecSnap, error)
	GetRecs(context.Context, tx.Conn, []identity.ADT) ([]DecRec, error)
	SelectEnv(context.Context, tx.Conn, []identity.ADT) (map[identity.ADT]DecRec, error)
}

