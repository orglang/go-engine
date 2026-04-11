package poolimpl

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
)

type Repo interface {
	GetRecByRef(db.Source, implsem.SemRef) (ImplRec, error)
}

type implJoinDS struct {
	ImplRef  implsem.SemRefDS `db:"sem"`
	LiabMode int16            `db:"exec.liab_mode"`
}

type implRecDS struct {
	ImplID     string             `db:"impl_id"`
	ImplRN     int64              `db:"impl_rn"`
	LiabMode   int16              `db:"liab_mode"`
	StructVars []implvar.VarRecDS `db:"struct_vars"`
	LinearVars []implvar.VarRecDS `db:"linear_vars"`
}
