package compexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/semterm"
)

type Repo interface {
	InsertRec(db.Source, ExecRec) error
	SelectSnap(db.Source, semterm.TermRef) (ExecSnap, error)
	UpdateProc(db.Source, ExecMod) error
}

type execRecDS struct {
	ImplID   string `db:"impl_id"`
	LiabMode int16  `db:"mode"`
}

type execModDS struct {
	ImplRefs   []semterm.TermRefDS
	LinearVars []compvar.VarRecDS
}
