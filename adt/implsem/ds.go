package implsem

import (
	"database/sql"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, SemRec) error
	TouchRec(db.Source, SemRef) error
	GetRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]SemRef, error)
}

type SemRefDS struct {
	ImplID string `db:"impl_id"`
	ImplRN int64  `db:"impl_rn"`
}

type semRecDS struct {
	ImplID string         `db:"impl_id"`
	ImplRN int64          `db:"impl_rn"`
	ImplQN sql.NullString `db:"impl_qn"`
	Kind   int16          `db:"kind"`
}

type SemBindDS struct {
	ImplID string         `db:"impl_id"`
	ImplQN sql.NullString `db:"impl_qn"`
}
