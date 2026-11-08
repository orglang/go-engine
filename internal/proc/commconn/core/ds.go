package core

import (
	"context"
	"database/sql"

	"orglang/go-engine/internal/lib/tx"

	commturn "orglang/go-engine/internal/proc/commturn/core"
)

type Repo interface {
	ModifyRec(context.Context, tx.Conn, ConnMod) error
	GetSnapByQry(context.Context, tx.Conn, ConnQry) (ConnSnap, error)
}

type ConnModDS struct {
	ConnID       string          `db:"conn_id"`
	OffsetNr     sql.Null[int64] `db:"offset_nr"`
	PrevOffsetNr int64           `db:"-"`
}

type ConnQryDS struct {
	ConnID string           `db:"conn_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type ConnSnapDS struct {
	ConnID   string               `db:"conn_id"`
	OffsetNr int64                `db:"offset_nr"`
	ActNr    int64                `db:"act_nr"`
	Turns    []commturn.TurnRecDS `db:"turns"`
}
