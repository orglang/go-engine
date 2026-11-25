package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/termsem"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DefRec) error
	SetExpCK(context.Context, tx.Conn, identity.ADT, contkey.ADT) error
	GetRefs(context.Context, tx.Conn) ([]termsem.SemRef, error)
	GetSnap(context.Context, tx.Conn, termsem.SemRef) (DefSnap, error)
	GetRecs(context.Context, tx.Conn, []identity.ADT) ([]DefRec, error)
	GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error)
	SelectEnv(context.Context, tx.Conn, []identity.ADT) (map[identity.ADT]DefRec, error)
}

type DefRepo interface {
	GetExpCKByID(context.Context, tx.Conn, identity.ADT) (contkey.ADT, error)
}
