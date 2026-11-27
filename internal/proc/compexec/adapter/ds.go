package adapter

import (
	comprootadapter "orglang/go-engine/internal/adt/comproot/adapter"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
)

type ExecRec struct {
	CompID   string `db:"comp_id"`
	CompRN   int64  `db:"comp_rn"`
	LiabMode int16  `db:"liab_mode"`
}

type ExecMod struct {
	CompRefs   []comprootadapter.SemRef
	LinearVars []compvaradapter.VarRec
}
