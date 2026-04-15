package termdec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/semtype"
	"orglang/go-engine/adt/termvar"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, DecRec) error
	GetRecByQN(db.Source, uniqsym.ADT) (DecRec, error)
}

type decRecDS struct {
	DescID    string             `db:"desc_id"`
	LiabVar   termvar.VarRecDS   `db:"liab_var"`
	AssetVars []termvar.VarRecDS `db:"asset_vars"`
}

type decSnapDS struct {
	DescRef   semtype.SemRefDS   `db:"ref"`
	LiabVar   termvar.VarRecDS   `db:"liab_var"`
	AssetVars []termvar.VarRecDS `db:"asset_vars"`
}
