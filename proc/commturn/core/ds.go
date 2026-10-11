package core

import (
	"database/sql"

	"orglang/go-engine/lib/db"

	termexp "orglang/go-engine/proc/termexp/core"
)

type Repo interface {
	InsertRecs(db.UoW, ...TurnRec) error
}

type StepRecDS struct {
	K      StepKindDS       `db:"kind"`
	ExecID sql.NullString   `db:"impl_id"`
	ChnlID sql.NullString   `db:"chnl_id"`
	ProcER termexp.ExpRecDS `db:"proc_er"`
}

type StepKindDS int

const (
	NonStep StepKindDS = iota
	MsgStep
	SvcStep
)
