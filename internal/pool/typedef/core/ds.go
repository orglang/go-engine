package core

import (
	"context"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DefRec) error
	GetRecsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

// TypeExpRepo stores and retrieves the AST belonging to a type definition.
type TypeExpRepo interface {
	AddRec(context.Context, tx.Conn, ExpRec, typesem.SemRef) error
	GetRecByCK(context.Context, tx.Conn, contkey.ADT) (ExpRec, error)
	GetRecsByCKs(context.Context, tx.Conn, []contkey.ADT) ([]ExpRec, error)
	GetRecMap(context.Context, tx.Conn, map[symbol.ADT]contkey.ADT) (map[symbol.ADT]ExpRec, error)
}
