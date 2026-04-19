package compexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, ExecRec) error
	ModifyRec(db.Source, ExecMod) error
	GetRecByRef(db.Source, compsem.SemRef) (ExecRec, error)
	GetSnapMapByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]ExecSnap1, error)
}

type execJoinDS struct {
	CompRef  compsem.SemRefDS `db:"ref"`
	LiabMode int16            `db:"exec.liab_mode"`
}

type execRecDS struct {
	CompID     string             `db:"comp_id"`
	CompRN     int64              `db:"comp_rn"`
	LiabMode   int16              `db:"liab_mode"`
	StructVars []compvar.VarRecDS `db:"struct_vars"`
	LinearVars []compvar.VarRecDS `db:"linear_vars"`
}

type execSnapDS struct {
	CompRef   compsem.SemRefDS  `db:"sem"`
	LiabMode  int16             `db:"exec.liab_mode"`
	StructVar *compvar.VarRecDS `db:"struct_var"`
	LinearVar *compvar.VarRecDS `db:"linear_var"`
}
