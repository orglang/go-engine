package adapter

import (
	comprootadapter "orglang/go-engine/internal/adt/comproot/adapter"
	execvaradapter "orglang/go-engine/internal/adt/execvar/adapter"
)

type EvalRec struct {
	CompID   string `db:"comp_id"`
	CompRN   int64  `db:"comp_rn"`
	TermQN   string `db:"term_qn"`
	LiabMode int16  `db:"liab_mode"`
}

type EvalSnap1 struct {
	CompRef   comprootadapter.RootRef `db:"sem"`
	LiabMode  int16                   `db:"exec.liab_mode"`
	StructVar *execvaradapter.VarRec  `db:"struct_var"`
	LinearVar *execvaradapter.VarRec  `db:"linear_var"`
}

type EvalSnap2 struct {
	CompID     string                  `db:"exec_id"`
	CompRN     int64                   `db:"exec_rn"`
	LiabMode   int16                   `db:"liab_mode"`
	StructVars []execvaradapter.VarRec `db:"struct_vars"`
	LinearVars []execvaradapter.VarRec `db:"linear_vars"`
}
