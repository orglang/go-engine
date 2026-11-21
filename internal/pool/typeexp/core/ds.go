package core

import (
	"context"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/typesem"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ExpRec, typesem.SemRef) error
	GetRecByCK(context.Context, tx.Conn, contkey.ADT) (ExpRec, error)
	GetRecsByCKs(context.Context, tx.Conn, []contkey.ADT) ([]ExpRec, error)
	GetRecMap(context.Context, tx.Conn, map[symbol.ADT]contkey.ADT) (map[symbol.ADT]ExpRec, error)
}
