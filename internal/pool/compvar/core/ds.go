package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/compvar"
)

type Repo interface {
	AddRecs(context.Context, tx.Conn, []compvar.VarRec) error
}
