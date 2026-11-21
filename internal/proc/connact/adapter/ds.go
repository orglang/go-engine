package adapter

import termexp "orglang/go-engine/internal/proc/termexp/core"

type ActRec struct {
	ConnID string           `db:"conn_id"`
	ActNr  int64            `db:"act_nr"`
	CompID string           `db:"comp_id"`
	CompRN int64            `db:"comp_rn"`
	ChnlID string           `db:"chnl_id"`
	K      StepKind         `db:"kind"`
	Exp    termexp.ExpRecDS `db:"exp" fieldopt:"noexpand"`
}

type StepKind int16

const (
	NonStep StepKind = iota
	MsgStep
	SvcStep
)
