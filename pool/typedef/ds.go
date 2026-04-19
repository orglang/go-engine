package typedef

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/typesem"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, DefRec) error
	Update(db.Source, DefRec) error
	SelectRefs(db.Source) ([]typesem.SemRef, error)
	SelectRecByRef(db.Source, typesem.SemRef) (DefRec, error)
	SelectRecsByRefs(db.Source, []typesem.SemRef) ([]DefRec, error)
	SelectRecByQN(db.Source, uniqsym.ADT) (DefRec, error)
	GetRecsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

type defRecDS struct {
	TypeID string `db:"type_id"`
	TypeRN string `db:"type_rn"`
	ExpVK  int64  `db:"exp_vk"`
}
