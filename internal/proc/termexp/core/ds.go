package core

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	// Insert persists an execution-oriented record. It is retained for the
	// runtime path and is not the Definition persistence contract.
	Insert(context.Context, tx.Conn, ExpRec) error

	// InsertSpec persists a declarative Definition tree and returns its root key.
	InsertSpec(context.Context, tx.Conn, DescSpec) error

	// SelectSpecByCK reads a declarative Definition tree by its root content key.
	SelectSpecByCK(context.Context, tx.Conn, contkey.ADT) (DescSpec, error)
}
