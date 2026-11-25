package core

import (
	"context"
	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DecRec) error
	GetRecByQN(context.Context, tx.Conn, uniqsym.ADT) (DecRec, error)
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
