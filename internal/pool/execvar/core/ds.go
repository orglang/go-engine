package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/execvar"
)

type Repo interface {
	AddRecs(context.Context, tx.Conn, []execvar.VarRec) error
}
