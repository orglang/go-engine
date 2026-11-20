package adapter

import termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"

type DecRecDS struct {
	TermID    string                    `db:"term_id"`
	TermRN    int64                     `db:"term_rn"`
	LiabVar   termvaradapter.VarRecDS   `db:"liab_var" fieldopt:"noexpand"`
	AssetVars []termvaradapter.VarRecDS `db:"asset_vars"`
}

type DecSnapDS struct {
	TermID    string                    `db:"term_id"`
	TermRN    int64                     `db:"term_rn"`
	LiabVar   termvaradapter.VarRecDS   `db:"liab_var"`
	AssetVars []termvaradapter.VarRecDS `db:"asset_vars"`
}
