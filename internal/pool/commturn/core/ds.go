package core

import termexp "orglang/go-engine/internal/pool/termexp/core"

type TurnRecDS struct {
	ConnID string           `db:"conn_id"`
	ActNr  int64            `db:"act_nr"`
	CompID string           `db:"comp_id"`
	CompRN int64            `db:"comp_rn"`
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
