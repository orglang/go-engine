package procexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
)

type Repo interface {
	InsertRec(db.Source, ExecRec) error
	SelectSnap(db.Source, implsem.SemRef) (ExecSnap, error)
	UpdateProc(db.Source, ExecMod) error
}

type execRecDS struct {
	ImplID   string `db:"impl_id"`
	LiabMode int16  `db:"mode"`
}

type execModDS struct {
	ImplRefs   []implsem.SemRefDS
	LinearVars []implvar.VarRecDS
}
