package adapter

import compexec "orglang/go-engine/internal/pool/compexec/adapter"

type ActRec struct {
	ConnID string         `db:"conn_id"`
	ActNr  int64          `db:"act_nr"`
	CompID string         `db:"comp_id"`
	CompRN int64          `db:"comp_rn"`
	ChnlID string         `db:"chnl_id"`
	K      actKind        `db:"kind"`
	Exp    compexec.ExpRec `db:"exp" fieldopt:"noexpand"`
}

type actKind int16

const (
	unkKind actKind = iota
	PubKind
	SubKind
)
