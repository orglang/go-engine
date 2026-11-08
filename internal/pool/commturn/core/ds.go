package core

import (
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

type TurnRecDS struct {
	ConnID string           `db:"conn_id"`
	ConnRN int64            `db:"conn_rn"`
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
