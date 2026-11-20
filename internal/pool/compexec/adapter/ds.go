package adapter

import (
	"orglang/go-engine/internal/adt/compsem"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
)

type ExecRecDS struct {
	CompID   string `db:"comp_id"`
	CompRN   int64  `db:"comp_rn"`
	LiabMode int16  `db:"liab_mode"`
}

type ExecSnap1DS struct {
	CompRef   compsem.SemRefDS  `db:"sem"`
	LiabMode  int16             `db:"exec.liab_mode"`
	StructVar *compvaradapter.VarRecDS `db:"struct_var"`
	LinearVar *compvaradapter.VarRecDS `db:"linear_var"`
}

type ExecSnap2DS struct {
	CompID     string             `db:"comp_id"`
	CompRN     int64              `db:"comp_rn"`
	LiabMode   int16              `db:"liab_mode"`
	StructVars []compvaradapter.VarRecDS `db:"struct_vars"`
	LinearVars []compvaradapter.VarRecDS `db:"linear_vars"`
}
