package adapter

import (
	execrootadapter "orglang/go-engine/internal/adt/execroot/adapter"
	execvaradapter "orglang/go-engine/internal/adt/execvar/adapter"
)

type EvalRec struct {
	ExecID   string `db:"exec_id"`
	ExecRN   int64  `db:"exec_rn"`
	LiabMode int16  `db:"liab_mode"`
}

type EvalMod struct {
	ExecRefs   []execrootadapter.RootRef
	LinearVars []execvaradapter.VarRec
}
