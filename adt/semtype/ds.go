package semtype

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, SemRec) error
	GetRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]TypeRef, error)
}

type SemRefDS struct {
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
}

type SemBindDS struct {
	DescQN string `db:"desc_qn"`
	DescID string `db:"desc_id"`
}

type semRecDS struct {
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
	DescQN string `db:"desc_qn"`
	Kind   uint16 `db:"kind"`
}
