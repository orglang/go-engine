package implsem

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	InsertRec(db.Source, SemRec) error
	SelectRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]SemRef, error)
}

type SemRefDS struct {
	ImplID string `db:"impl_id"`
	ImplRN int64  `db:"impl_rn"`
}

type semRecDS struct {
	ImplID string `db:"impl_id"`
	ImplRN int64  `db:"impl_rn"`
	Kind   uint8  `db:"kind"`
}
