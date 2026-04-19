package typesem

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	TouchRef(db.Source, SemRef) error
	GetRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]SemRef, error)
}

type SemRefDS struct {
	TypeID string `db:"type_id"`
	TypeRN int64  `db:"type_rn"`
}
