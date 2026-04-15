package typedef

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/semtype"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, DefRec) error
	Update(db.Source, DefRec) error
	SelectRefs(db.Source) ([]semtype.TypeRef, error)
	SelectRecByRef(db.Source, semtype.TypeRef) (DefRec, error)
	SelectRecsByRefs(db.Source, []semtype.TypeRef) ([]DefRec, error)
	SelectRecByQN(db.Source, uniqsym.ADT) (DefRec, error)
	GetRecsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]DefRec, error)
}

type defRecDS struct {
	DescID string `db:"desc_id"`
	// DescRN int64  `db:"desc_rn"`
	ExpVK int64 `db:"exp_vk"`
}
