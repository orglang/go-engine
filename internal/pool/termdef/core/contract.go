package core

import (
	"context"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	GetByTermID(context.Context, tx.Conn, identity.ADT) (DefRec, error)
}
