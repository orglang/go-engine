package typedef

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	Insert(db.Source, DefRec) error
	Update(db.Source, DefRec) error
	SelectRefs(db.Source) ([]DefRef, error)
	SelectRecByRef(db.Source, DefRef) (DefRec, error)
	SelectRecsByRefs(db.Source, []DefRef) ([]DefRec, error)
	SelectRecByQN(db.Source, uniqsym.ADT) (DefRec, error)
	SelectRecsByQNs(db.Source, []uniqsym.ADT) ([]DefRec, error)
	SelectEnv(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

type defRefDS struct {
	ID string `db:"def_id"`
	RN int64  `db:"def_rn"`
}

type defRecDS struct {
	ID    string `db:"def_id"`
	RN    int64  `db:"def_rn"`
	ExpID string `db:"exp_id"`
	Title string `db:"title"`
}
