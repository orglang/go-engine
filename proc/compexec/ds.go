package compexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/compvar"
)

type Repo interface {
	AddRec(db.Source, ExecRec) error
	SelectSnap(db.Source, compsem.SemRef) (ExecSnap, error)
	UpdateProc(db.Source, ExecMod) error
}

type execRecDS struct {
	CompID   string `db:"impl_id"`
	LiabMode int16  `db:"mode"`
}

type execModDS struct {
	CompRefs   []compsem.SemRefDS
	LinearVars []compvar.VarRecDS
}
