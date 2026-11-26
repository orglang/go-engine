package core

import (
	"context"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	SetExpCK(context.Context, tx.Conn, identity.ADT, contkey.ADT) error
	GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error)
	InsertSpec(context.Context, tx.Conn, ExpSpec) error
	SelectSpecByCK(context.Context, tx.Conn, contkey.ADT) (ExpSpec, error)
}
