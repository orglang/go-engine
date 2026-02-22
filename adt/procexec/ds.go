package procexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/procstep"
)

type Repo interface {
	InsertRec(db.Source, ExecRec) error
	SelectSnap(db.Source, implsem.SemRef) (ExecSnap, error)
	UpdateProc(db.Source, ExecMod) error
}

type execRecDS struct {
}

type execModDS struct {
	Locks []implsem.SemRefDS
	Binds []implvar.VarRecDS
	Steps []procstep.StepRecDS
}

type liabDS struct {
	PoolID string `db:"desc_id"`
	ProcID string `db:"proc_id"`
	PoolRN int64  `db:"rev"`
}
