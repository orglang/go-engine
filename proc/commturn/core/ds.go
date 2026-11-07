package core

import (
	"context"
	"database/sql"

	"orglang/go-engine/lib/tx"

	termexp "orglang/go-engine/proc/termexp/core"
)

type Repo interface {
	InsertRecs(context.Context, tx.Conn, ...TurnRec) error
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
