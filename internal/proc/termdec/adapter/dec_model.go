package adapter

import termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"

type DecRec struct {
	TermID    string                  `db:"term_id"`
	ProgID    string                  `db:"prog_id"`
	ProgRN    int64                   `db:"prog_rn"`
	LiabVar   termvaradapter.VarRec   `db:"liab_var" fieldopt:"noexpand"`
	AssetVars []termvaradapter.VarRec `db:"asset_vars"`
}

type DecSnap struct {
	TermID    string                  `db:"term_id"`
	ProgID    string                  `db:"prog_id"`
	ProgRN    int64                   `db:"prog_rn"`
	LiabVar   termvaradapter.VarRec   `db:"liab_var"`
	AssetVars []termvaradapter.VarRec `db:"asset_vars"`
}
