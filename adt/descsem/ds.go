package descsem

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	InsertRec(db.Source, SemRec) error
	SelectRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]SemRef, error)
}

type SemRefDS struct {
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
}

type semBindDS struct {
	DescQN string `db:"desc_qn"`
	DescID string `db:"desc_id"`
}

type semRecDS struct {
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
	DescQN string `db:"desc_qn"`
	Kind   uint8  `db:"kind"`
}
