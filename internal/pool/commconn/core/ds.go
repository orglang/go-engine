package core

import (
	"context"
	"database/sql"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/connsem"
	"orglang/go-engine/internal/adt/uniqsym"
	commturn "orglang/go-engine/internal/pool/commturn/core"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ConnRec) error
	ModifyRec(context.Context, tx.Conn, ConnMod) error
	GetRefsByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]connsem.SemRef, error)
	GetSnapByQry(context.Context, tx.Conn, ConnQry) (ConnSnap, error)
}

type ConnRecDS struct {
	ConnID   string `db:"comm_id"`
	ConnRN   int64  `db:"comm_rn"`
	OffsetNr int64  `db:"offset_nr"`
}

type ConnModDS struct {
	ConnID   string          `db:"comm_id"`
	OffsetNr sql.Null[int64] `db:"offset_nr"`
}

type ConnQryDS struct {
	ConnID string           `db:"comm_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type ConnSnapDS struct {
	ConnID   string               `db:"comm_id"`
	ConnRN   int64                `db:"comm_rn"`
	OffsetNr int64                `db:"offset_nr"`
	Turns    []commturn.TurnRecDS `db:"turns"`
}
