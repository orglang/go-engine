package descexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	InsertRec(db.Source, ExecRec) error
	SelectRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]ExecRef, error)
}

type ExecRefDS struct {
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
}

type execRecDS struct {
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
	Kind   uint8  `db:"kind"`
}
