package core

import (
	"context"
	"database/sql"

	"orglang/go-engine/internal/lib/tx"

	commturn "orglang/go-engine/internal/proc/commturn/core"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ConnRec) error
	ModifyRec(context.Context, tx.Conn, ConnMod) error
	GetSnapByQry(context.Context, tx.Conn, ConnQry) (ConnSnap, error)
}

type ConnRecDS struct {
	CommID   string `db:"comm_id"`
	CommRN   int64  `db:"comm_rn"`
	OffsetNr int64  `db:"offset_nr"`
}

type ConnModDS struct {
	CommID   string          `db:"comm_id"`
	OffsetNr sql.Null[int64] `db:"offset_nr"`
}

type ConnQryDS struct {
	CommID string           `db:"comm_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type ConnSnapDS struct {
	CommID   string                `db:"comm_id"`
	CommRN   int64                 `db:"comm_rn"`
	OffsetNr int64                 `db:"offset_nr"`
	Turns    []commturn.TurnRecDS  `db:"turns"`
}
