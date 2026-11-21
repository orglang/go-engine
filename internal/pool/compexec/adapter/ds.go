package adapter

import (
	compsemadapter "orglang/go-engine/internal/adt/compsem/adapter"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
)

type ExecRec struct {
	CompID   string `db:"comp_id"`
	CompRN   int64  `db:"comp_rn"`
	LiabMode int16  `db:"liab_mode"`
}

type ExecSnap1 struct {
	CompRef   compsemadapter.SemRef  `db:"sem"`
	LiabMode  int16                    `db:"exec.liab_mode"`
	StructVar *compvaradapter.VarRec `db:"struct_var"`
	LinearVar *compvaradapter.VarRec `db:"linear_var"`
}

type ExecSnap2 struct {
	CompID     string                    `db:"comp_id"`
	CompRN     int64                     `db:"comp_rn"`
	LiabMode   int16                     `db:"liab_mode"`
	StructVars []compvaradapter.VarRec `db:"struct_vars"`
	LinearVars []compvaradapter.VarRec `db:"linear_vars"`
}
