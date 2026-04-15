package poolexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/semterm"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, ExecRec) error
	GetRecMapByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]ExecRec, error)
	GetRefs(db.Source) ([]semterm.TermRef, error)
	GetSnapMapByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]ExecSnap, error)
}

type execRecDS struct {
	ImplID   string `db:"impl_id"`
	LiabMode int16  `db:"liab_mode"`
}

type execSnapDS struct {
	ImplRef   semterm.TermRefDS `db:"sem"`
	LiabMode  int16             `db:"exec.liab_mode"`
	StructVar *compvar.VarRecDS `db:"struct_var"`
	LinearVar *compvar.VarRecDS `db:"linear_var"`
}
