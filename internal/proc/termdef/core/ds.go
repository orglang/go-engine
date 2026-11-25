package core

import (
	"context"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DecRec) error
	GetRefs(context.Context, tx.Conn) ([]termsem.SemRef, error)
	GetSnap(context.Context, tx.Conn, termsem.SemRef) (DecSnap, error)
	GetRecs(context.Context, tx.Conn, []identity.ADT) ([]DecRec, error)
	SelectEnv(context.Context, tx.Conn, []identity.ADT) (map[identity.ADT]DecRec, error)

	SetExpCK(context.Context, tx.Conn, identity.ADT, contkey.ADT) error
	GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error)
	GetQNByID(context.Context, tx.Conn, identity.ADT) (uniqsym.ADT, error)
	InsertSpec(context.Context, tx.Conn, ExpSpec) error
	SelectSpecByCK(context.Context, tx.Conn, contkey.ADT) (ExpSpec, error)
}
