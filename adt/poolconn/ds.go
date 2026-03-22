package poolconn

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, ConnRec) error
	UpdateRec(db.Source, ConnMod) error
	GetRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error)
	GetSnapByQry(db.Source, ConnQry) (ConnSnap, error)
}

type connRecDS struct {
	CommID string `db:"comm_id"`
	CommRN int64  `db:"comm_on"`
}
