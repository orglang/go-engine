package core

import (
	"context"

	"orglang/go-engine/internal/lib/tx"
)

// Repo persists and retrieves declarative pool term-definition trees.
type Repo interface {
	InsertSpec(context.Context, tx.Conn, DescSpec) error
	SelectSpecByCK(context.Context, tx.Conn, contkey.ADT) (DescSpec, error)
}
