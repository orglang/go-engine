package core

import (
	"context"
	"orglang/go-engine/internal/adt/contkey"
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

type DefRepo interface {
	SetExpCK(context.Context, tx.Conn, identity.ADT, contkey.ADT) error
	GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error)
}

// ExpRepo stores and retrieves recursive Definition nodes.
type ExpRepo interface {
	InsertSpec(context.Context, tx.Conn, ExpSpec) error
	SelectSpecByCK(context.Context, tx.Conn, contkey.ADT) (ExpSpec, error)
}
