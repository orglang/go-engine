package core

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/termsem"
	"orglang/go-engine/adt/termvar"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.UoW, DefRec) error
	GetRecByQN(db.UoW, uniqsym.ADT) (DefRec, error)
}

type DefRecDS struct {
	TermID    string             `db:"term_id"`
	TermRN    int64              `db:"term_rn"`
	LiabVar   termvar.VarRecDS   `db:"liab_var" fieldopt:"noexpand"`
	AssetVars []termvar.VarRecDS `db:"asset_vars"`
}

type DefSnapDS struct {
	TermRef   termsem.SemRefDS   `db:"ref"`
	LiabVar   termvar.VarRecDS   `db:"liab_var"`
	AssetVars []termvar.VarRecDS `db:"asset_vars"`
}
