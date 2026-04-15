package commexch

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/semcomm"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	InsertRec(db.Source, ExchRec) error
	UpdateRec(db.Source, ExchMod) error
	SelectRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]semcomm.CommRef, error)
	SelectSnapByQry(db.Source, ExchQry) (ExchSnap, error)
}
