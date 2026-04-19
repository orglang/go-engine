package typedef

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/typesem"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, DefRec) error
	ModifyRec(db.Source, DefRec) error
	GetRefs(db.Source) ([]typesem.SemRef, error)
	GetRecByRef(db.Source, typesem.SemRef) (DefRec, error)
	GetRecsByRefs(db.Source, []typesem.SemRef) ([]DefRec, error)
	GetRecByQN(db.Source, uniqsym.ADT) (DefRec, error)
	GetRecsByQNs(db.Source, []uniqsym.ADT) ([]DefRec, error)
	SelectEnv(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

type defRecDS struct {
	TypeID string `db:"type_id"`
	TypeRN int64  `db:"type_rn"`
	ExpVK  int64  `db:"exp_vk"`
}
