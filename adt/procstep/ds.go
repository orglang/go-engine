package procstep

import (
	"database/sql"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/procexp"
)

type Repo interface {
	InsertRecs(db.Source, ...StepRec) error
}

type StepRecDS struct {
	K      stepKindDS       `db:"kind"`
	ExecID sql.NullString   `db:"exec_id"`
	ChnlID sql.NullString   `db:"chnl_id"`
	ProcER procexp.ExpRecDS `db:"proc_er"`
}

type stepKindDS int

const (
	nonStep = stepKindDS(iota)
	msgStep
	svcStep
)
