package adapter

import execeval "orglang/go-engine/internal/pool/execeval/adapter"

type ActRec struct {
	ConnID string          `db:"conn_id"`
	ActNr  int64           `db:"act_nr"`
	ExecID string          `db:"exec_id"`
	ExecRN int64           `db:"exec_rn"`
	ChnlID string          `db:"chnl_id"`
	K      actKind         `db:"kind"`
	Exp    execeval.ExpRec `db:"exp" fieldopt:"noexpand"`
}

type actKind int16

const (
	unkKind actKind = iota
	PubKind
	SubKind
)
