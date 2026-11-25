package adapter

import termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"

type DecRec struct {
	TermID    string                  `db:"term_id"`
	DescID    string                  `db:"desc_id"`
	DescRN    int64                   `db:"desc_rn"`
	ExpCK     int64                   `db:"exp_ck"`
	LiabVar   termvaradapter.VarRec   `db:"liab_var" fieldopt:"noexpand"`
	AssetVars []termvaradapter.VarRec `db:"asset_vars"`
}
