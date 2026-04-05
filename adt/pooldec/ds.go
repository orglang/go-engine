package pooldec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/descvar"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, DecRec) error
	GetRecByQN(db.Source, uniqsym.ADT) (DecRec, error)
}

type decRecDS struct {
	DescID    string             `db:"desc_id"`
	LiabVar   descvar.VarRecDS   `db:"liab_var"`
	AssetVars []descvar.VarRecDS `db:"asset_vars"`
}

type decSnapDS struct {
	DescRef   descsem.SemRefDS   `db:"ref"`
	LiabVar   descvar.VarRecDS   `db:"liab_var"`
	AssetVars []descvar.VarRecDS `db:"asset_vars"`
}
