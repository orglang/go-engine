package compexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/compvar"
)

type Repo interface {
	AddRec(db.Source, ExecRec) error
	ModifyRec(db.Source, ExecMod) error
	GetSnapByRef(db.Source, compsem.SemRef) (ExecSnap, error)
}

type execRecDS struct {
	CompID   string `db:"comp_id"`
	CompRN   int64  `db:"comp_rn"`
	LiabMode int16  `db:"liab_mode"`
}

type execModDS struct {
	CompRefs   []compsem.SemRefDS
	LinearVars []compvar.VarRecDS
}
