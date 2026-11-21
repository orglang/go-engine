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

type ExecMod struct {
	CompRefs   []compsemadapter.SemRef
	LinearVars []compvaradapter.VarRec
}
