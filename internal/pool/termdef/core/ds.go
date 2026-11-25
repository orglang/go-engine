package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DefRec) error
	SetExpCK(context.Context, tx.Conn, identity.ADT, contkey.ADT) error
	GetRecByQN(context.Context, tx.Conn, uniqsym.ADT) (DefRec, error)
	GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error)
}

type DefRepo interface {
	GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error)
}
