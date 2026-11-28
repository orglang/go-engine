package adapter

import (
	execrootadapter "orglang/go-engine/internal/adt/execroot/adapter"
	execvaradapter "orglang/go-engine/internal/adt/execvar/adapter"
)

type EvalRec struct {
	RootID   string `db:"comp_id"`
	RootRN   int64  `db:"comp_rn"`
	LiabMode int16  `db:"liab_mode"`
}

type EvalMod struct {
	RootRefs   []execrootadapter.RootRef
	LinearVars []execvaradapter.VarRec
}
