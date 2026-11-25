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
	GetSnap(context.Context, tx.Conn, termsem.SemRef) (DefSnap, error)
	GetRecs(context.Context, tx.Conn, []identity.ADT) ([]DecRec, error)
	SelectEnv(context.Context, tx.Conn, []identity.ADT) (map[identity.ADT]DecRec, error)
}

type DefRepo interface {
	SetExpCK(context.Context, tx.Conn, identity.ADT, contkey.ADT) error
	GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error)
}
