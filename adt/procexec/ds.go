package procexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/procbind"
	"orglang/go-engine/adt/procstep"
	"orglang/go-engine/adt/uniqref"
)

type Repo interface {
	SelectSnap(db.Source, ExecRef) (ExecSnap, error)
	UpdateProc(db.Source, ExecMod) error
}

type execModDS struct {
	Locks []execRefDS
	Binds []procbind.BindRecDS
	Steps []procstep.StepRecDS
}

type execRefDS = uniqref.Data

type liabDS struct {
	PoolID string `db:"pool_id"`
	ProcID string `db:"proc_id"`
	PoolRN int64  `db:"rev"`
}
