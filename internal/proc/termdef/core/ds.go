package core

import (
	"context"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	SetExprCK(context.Context, tx.Conn, identity.ADT, contkey.ADT) error
	GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error)
	InsertSpec(context.Context, tx.Conn, ExprSpec) error
	SelectSpecByCK(context.Context, tx.Conn, contkey.ADT) (ExprSpec, error)
}
