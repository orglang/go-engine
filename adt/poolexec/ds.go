package poolexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, ExecRec) error
	GetRecMapByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]ExecRec, error)
	GetRefs(db.Source) ([]implsem.SemRef, error)
	GetSnapMapByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]ExecSnap, error)
}

type execRecDS struct {
	ImplID   string `db:"impl_id"`
	LiabMode int16  `db:"liab_mode"`
}

type execSnapDS struct {
	ImplRef   implsem.SemRefDS  `db:"sem"`
	LiabMode  int16             `db:"exec.liab_mode"`
	StructVar *implvar.VarRecDS `db:"struct_var"`
	LinearVar *implvar.VarRecDS `db:"linear_var"`
}
