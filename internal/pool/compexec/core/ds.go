package core

import (
	"context"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/compvar"
	"orglang/go-engine/internal/adt/uniqsym"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ExecRec) error
	ModifyRec(context.Context, tx.Conn, ExecMod) error
	GetSnapByRef(context.Context, tx.Conn, compsem.SemRef) (ExecSnap2, error)
	GetSnapMapByQNs(context.Context, tx.Conn, []uniqsym.ADT) (map[uniqsym.ADT]ExecSnap1, error)
}

type ExecRecDS struct {
	CompID   string `db:"comp_id"`
	CompRN   int64  `db:"comp_rn"`
	LiabMode int16  `db:"liab_mode"`
}

type ExecSnap1DS struct {
	CompRef   compsem.SemRefDS  `db:"sem"`
	LiabMode  int16             `db:"exec.liab_mode"`
	StructVar *compvar.VarRecDS `db:"struct_var"`
	LinearVar *compvar.VarRecDS `db:"linear_var"`
}

type ExecSnap2DS struct {
	CompID     string             `db:"comp_id"`
	CompRN     int64              `db:"comp_rn"`
	LiabMode   int16              `db:"liab_mode"`
	StructVars []compvar.VarRecDS `db:"struct_vars"`
	LinearVars []compvar.VarRecDS `db:"linear_vars"`
}
