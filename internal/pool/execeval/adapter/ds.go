package adapter

import (
	execrootadapter "orglang/go-engine/internal/adt/execroot/adapter"
	execvaradapter "orglang/go-engine/internal/adt/execvar/adapter"
)

type EvalRec struct {
	RootID   string `db:"root_id"`
	RootRN   int64  `db:"root_rn"`
	TermQN   string `db:"term_qn"`
	LiabMode int16  `db:"liab_mode"`
}

type EvalSnap1 struct {
	RootRef   execrootadapter.RootRef `db:"sem"`
	LiabMode  int16                   `db:"exec.liab_mode"`
	StructVar *execvaradapter.VarRec  `db:"struct_var"`
	LinearVar *execvaradapter.VarRec  `db:"linear_var"`
}

type EvalSnap2 struct {
	RootID     string                  `db:"exec_id"`
	RootRN     int64                   `db:"exec_rn"`
	LiabMode   int16                   `db:"liab_mode"`
	StructVars []execvaradapter.VarRec `db:"struct_vars"`
	LinearVars []execvaradapter.VarRec `db:"linear_vars"`
}
