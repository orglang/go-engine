package commsem

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	InsertRec(db.Source, SemRec) error
	TouchRec(db.Source, SemRef) error
	SelectRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]SemRef, error)
}

type SemRefDS struct {
	CommID string `db:"comm_id"`
	CommRN int64  `db:"comm_rn"`
}

type semRecDS struct {
	CommID string `db:"comm_id"`
	CommRN int64  `db:"comm_rn"`
	Kind   int8   `db:"kind"`
}
