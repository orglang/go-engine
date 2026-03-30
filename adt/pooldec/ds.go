package pooldec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descvar"
)

type Repo interface {
	AddRec(db.Source, DecRec) error
}

type decRecDS struct {
	DescID    string             `db:"desc_id"`
	DescRN    int64              `db:"desc_rn"`
	LiabVar   descvar.VarRecDS   `db:"liab_var"`
	AssetVars []descvar.VarRecDS `db:"asset_vars"`
}
