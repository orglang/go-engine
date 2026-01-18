package typedef

import (
	"orglang/go-runtime/lib/db"

	"orglang/go-runtime/adt/uniqref"
	"orglang/go-runtime/adt/uniqsym"
)

type Repo interface {
	Insert(db.Source, DefRec) error
	Update(db.Source, DefRec) error
	SelectRefs(db.Source) ([]DefRef, error)
	SelectRecByID(db.Source, DefRef) (DefRec, error)
	SelectRecsByIDs(db.Source, []DefRef) ([]DefRec, error)
	SelectRecByQN(db.Source, uniqsym.ADT) (DefRec, error)
	SelectRecsByQNs(db.Source, []uniqsym.ADT) ([]DefRec, error)
	SelectEnv(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

type defRefDS = uniqref.Data

type defRecDS struct {
	DefID string `db:"def_id"`
	Title string `db:"title"`
	ExpID string `db:"exp_id"`
	DefRN int64  `db:"def_rn"`
}
