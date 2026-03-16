package procconn

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	InsertRec(db.Source, ConnRec) error
	UpdateRec(db.Source, ConnMod) error
	SelectRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error)
	SelectSnapByQry(db.Source, ConnQuery) (ConnSnap, error)
}
