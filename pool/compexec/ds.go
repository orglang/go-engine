package compexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/semcomp"
	"orglang/go-engine/adt/semterm"
)

type Repo interface {
	ModifyRec(db.Source, ExecMod) error
	GetRecByRef(db.Source, semcomp.CompRef) (ExecRec, error)
}

type implJoinDS struct {
	ImplRef  semterm.TermRefDS `db:"sem"`
	LiabMode int16             `db:"exec.liab_mode"`
}

type implRecDS struct {
	ImplID     string             `db:"impl_id"`
	ImplRN     int64              `db:"impl_rn"`
	LiabMode   int16              `db:"liab_mode"`
	StructVars []compvar.VarRecDS `db:"struct_vars"`
	LinearVars []compvar.VarRecDS `db:"linear_vars"`
}
