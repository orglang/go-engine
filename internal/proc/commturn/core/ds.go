package core

import termexp "orglang/go-engine/internal/proc/termexp/core"

type TurnRecDS struct {
	ConnID string           `db:"comm_id"`
	ConnRN int64            `db:"comm_rn"`
	CompID string           `db:"comp_id"`
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
