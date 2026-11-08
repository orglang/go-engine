package core

import (
	"context"
	"database/sql"
	"orglang/go-engine/lib/tx"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/uniqsym"
	commturn "orglang/go-engine/pool/commturn/core"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ExchRec) error
	ModifyRec(context.Context, tx.Conn, ExchMod) error
	GetRefsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error)
	GetSnapByQry(context.Context, tx.Conn, ExchQry) (ExchSnap, error)
}

type ExchRecDS struct {
	CommID   string `db:"comm_id"`
	CommRN   int64  `db:"comm_rn"`
	OffsetNr int64  `db:"offset_nr"`
}

type ExchModDS struct {
	CommID   string          `db:"comm_id"`
	OffsetNr sql.Null[int64] `db:"offset_nr"`
}

type ExchQryDS struct {
	CommID string           `db:"comm_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type ExchSnapDS struct {
	CommID   string               `db:"comm_id"`
	CommRN   int64                `db:"comm_rn"`
	OffsetNr int64                `db:"offset_nr"`
	Turns    []commturn.TurnRecDS `db:"turns"`
}
