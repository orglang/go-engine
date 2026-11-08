package core

import (
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

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
