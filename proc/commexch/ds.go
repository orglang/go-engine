package commexch

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, ExchRec) error
	Modifyec(db.Source, ExchMod) error
	GetRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error)
	GetSnapByQry(db.Source, ExchQry) (ExchSnap, error)
}
