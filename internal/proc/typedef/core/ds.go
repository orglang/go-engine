package core

import (
	"context"

	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
)

type DescSemRepo interface {
	AddRec(context.Context, tx.Conn, descsem.SemRec) error
}

type Repo interface {
	AddRec(context.Context, tx.Conn, DefRec) error
	ModifyRec(context.Context, tx.Conn, DefRec) error
	GetRefs(context.Context, tx.Conn) ([]typesem.SemRef, error)
	GetRecByRef(context.Context, tx.Conn, typesem.SemRef) (DefRec, error)
	GetRecsByRefs(context.Context, tx.Conn, []typesem.SemRef) ([]DefRec, error)
	GetRecByQN(context.Context, tx.Conn, uniqsym.ADT) (DefRec, error)
	GetRecsByQNs(context.Context, tx.Conn, []uniqsym.ADT) ([]DefRec, error)
	SelectEnv(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

