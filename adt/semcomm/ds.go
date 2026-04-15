package semcomm

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, SemRec) error
	TouchRec(db.Source, CommRef) error
	SelectRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]CommRef, error)
}

type SemRefDS struct {
	CommID string `db:"comm_id"`
	CommRN int64  `db:"comm_rn"`
}

type semRecDS struct {
	CommID string `db:"comm_id"`
	CommRN int64  `db:"comm_rn"`
	Kind   int16  `db:"kind"`
}
