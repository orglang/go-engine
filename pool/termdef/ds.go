package termdef

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/termvar"
	"orglang/go-engine/adt/typesem"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, DefRec) error
	GetRecByQN(db.Source, uniqsym.ADT) (DefRec, error)
}

type decRecDS struct {
	TermID    string             `db:"desc_id"`
	LiabVar   termvar.VarRecDS   `db:"liab_var"`
	AssetVars []termvar.VarRecDS `db:"asset_vars"`
}

type decSnapDS struct {
	TermRef   typesem.SemRefDS   `db:"ref"`
	LiabVar   termvar.VarRecDS   `db:"liab_var"`
	AssetVars []termvar.VarRecDS `db:"asset_vars"`
}
