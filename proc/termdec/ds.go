package termdec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/semtype"
	"orglang/go-engine/adt/termvar"
)

type Repo interface {
	InsertRec(db.Source, DecRec) error
	SelectRefs(db.Source) ([]semtype.TypeRef, error)
	SelectSnap(db.Source, semtype.TypeRef) (DecSnap, error)
	SelectRecs(db.Source, []identity.ADT) ([]DecRec, error)
	SelectEnv(db.Source, []identity.ADT) (map[identity.ADT]DecRec, error)
}

type decRecDS struct {
	DescID    string             `db:"desc_id"`
	DescRN    int64              `db:"desc_rn"`
	LiabVar   termvar.VarRecDS   `db:"liab_var"`
	AssetVars []termvar.VarRecDS `db:"asset_vars"`
}

type decSnapDS struct {
	DescID    string             `db:"desc_id"`
	DescRN    int64              `db:"desc_rn"`
	LiabVar   termvar.VarRecDS   `db:"liab_var"`
	AssetVars []termvar.VarRecDS `db:"asset_vars"`
}
