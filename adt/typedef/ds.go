package typedef

import (
	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
)

type Repo interface {
	Insert(sd.Source, DefRec) error
	Update(sd.Source, DefRec) error
	SelectRefs(sd.Source) ([]DefRef, error)
	SelectRecByID(sd.Source, identity.ADT) (DefRec, error)
	SelectRecsByIDs(sd.Source, []identity.ADT) ([]DefRec, error)
	SelectRecByQN(sd.Source, qualsym.ADT) (DefRec, error)
	SelectRecsByQNs(sd.Source, []qualsym.ADT) ([]DefRec, error)
	SelectEnv(sd.Source, []qualsym.ADT) (map[qualsym.ADT]DefRec, error)
}

type defRefDS struct {
	DefID string `db:"def_id"`
	DefRN int64  `db:"def_rn"`
}

type defRecDS struct {
	DefID string `db:"def_id"`
	Title string `db:"title"`
	ExpID string `db:"exp_id"`
	DefRN int64  `db:"def_rn"`
}
