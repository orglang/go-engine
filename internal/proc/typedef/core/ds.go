package core

import (
	"context"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, DefRec) error
	ModifyRec(context.Context, tx.Conn, DefRec) error
	GetIDs(context.Context, tx.Conn) ([]identity.ADT, error)
	GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error)
	GetRecsByIDs(context.Context, tx.Conn, []identity.ADT) ([]DefRec, error)
	GetRecByQN(context.Context, tx.Conn, uniqsym.ADT) (DefRec, error)
	GetRecsByQNs(context.Context, tx.Conn, []uniqsym.ADT) ([]DefRec, error)
	SelectEnv(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

type TypeExpRepo interface {
	AddRec(context.Context, tx.Conn, ExpRec) error
	SelectRecByCK(context.Context, tx.Conn, contkey.ADT) (ExpRec, error)
	SelectRecsByCKs(context.Context, tx.Conn, []contkey.ADT) ([]ExpRec, error)
	SelectEnv(context.Context, tx.Conn, []contkey.ADT) (map[contkey.ADT]ExpRec, error)
}
