package core

import (
	"context"
	"orglang/go-engine/lib/tx"

	termexp "orglang/go-engine/pool/termexp/core"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, TurnRec) error
	AddRecs(context.Context, tx.Conn, []TurnRec) error
}

type TurnRecDS struct {
	CommID string           `db:"comm_id"`
	CommRN int64            `db:"comm_rn"`
	CompID string           `db:"comp_id"`
	ChnlID string           `db:"chnl_id"`
	K      turnKind         `db:"kind"`
	Exp    termexp.ExpRecDS `db:"exp" fieldopt:"noexpand"`
}

type turnKind int16

const (
	unkKind turnKind = iota
	PubKind
	SubKind
)
