package core

import (
	"context"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ExpRec) error
	SelectRecByCK(context.Context, tx.Conn, contkey.ADT) (ExpRec, error)
	SelectRecsByCKs(context.Context, tx.Conn, []contkey.ADT) ([]ExpRec, error)
	SelectEnv(context.Context, tx.Conn, []contkey.ADT) (map[contkey.ADT]ExpRec, error)
}
