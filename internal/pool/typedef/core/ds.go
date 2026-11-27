package core

import (
	"context"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DefRec) error
	GetRecsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
	GetTypeMap(context.Context, tx.Conn, map[symbol.ADT]contkey.ADT) (map[symbol.ADT]ExpRec, error)
}
