package poolexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, ExecRec) error
	GetRecsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]ExecRec, error)
	GetRefs(db.Source) ([]implsem.SemRef, error)
	GetSnap(db.Source, implsem.SemRef) (ExecSnap, error)
	GetSnapsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]ExecLiabSnap, error)
}

type execRecDS struct {
	ImplID string `db:"impl_id"`
	Mode   int8   `db:"mode"`
}

type execSnapDS struct {
	ImplID     string             `db:"impl_id"`
	ImplRN     int64              `db:"impl_rn"`
	StructVars []implvar.VarRecDS `db:"struct_vars"`
	LinearVars []implvar.VarRecDS `db:"linear_vars"`
}

type execLiabSnapDS struct {
	ImplRef   implsem.SemRefDS  `db:"sem"`
	Mode      int8              `db:"exec.mode"`
	StructVar *implvar.VarRecDS `db:"struct_var"`
	LinearVar *implvar.VarRecDS `db:"linear_var"`
}
