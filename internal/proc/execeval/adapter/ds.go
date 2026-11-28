package adapter

import (
	execrootadapter "orglang/go-engine/internal/adt/execroot/adapter"
	execvaradapter "orglang/go-engine/internal/adt/execvar/adapter"
)

type EvalRec struct {
	RootID   string `db:"exec_id"`
	RootRN   int64  `db:"exec_rn"`
	LiabMode int16  `db:"liab_mode"`
}

type EvalMod struct {
	RootRefs   []execrootadapter.RootRef
	LinearVars []execvaradapter.VarRec
}
