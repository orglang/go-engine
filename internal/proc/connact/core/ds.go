package core

import termexp "orglang/go-engine/internal/proc/termexp/core"

type ActRecDS struct {
	ConnID string           `db:"conn_id"`
	ActNr  int64            `db:"act_nr"`
	CompID string           `db:"comp_id"`
	CompRN int64            `db:"comp_rn"`
	ChnlID string           `db:"chnl_id"`
	K      StepKindDS       `db:"kind"`
	Exp    termexp.ExpRecDS `db:"exp" fieldopt:"noexpand"`
}

type StepKindDS int16

const (
	NonStep StepKindDS = iota
	MsgStep
	SvcStep
)
