package core

import (
	"context"
	"orglang/go-engine/lib/tx"

	"orglang/go-engine/adt/compvar"
)

type Repo interface {
	AddRecs(context.Context, tx.Conn, []compvar.VarRec) error
}
