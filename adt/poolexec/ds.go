package poolexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	InsertRec(db.Source, ExecRec) error
	SelectRecsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]ExecRec, error)
	SelectRefs(db.Source) ([]implsem.SemRef, error)
	SelectSnap(db.Source, implsem.SemRef) (ExecSnap, error)
}

type execRecDS struct {
	ImplID string `db:"impl_id"`
}

type execSnapDS struct {
	ImplID     string             `db:"impl_id"`
	ImplRN     int64              `db:"impl_rn"`
	StructVars []implvar.VarRecDS `db:"struct_vars"`
	LinearVars []implvar.VarRecDS `db:"linear_vars"`
}
