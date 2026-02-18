package xactdef

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/descexec"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	InsertRec(db.Source, DefRec) error
	Update(db.Source, DefRec) error
	SelectRefs(db.Source) ([]descexec.ExecRef, error)
	SelectRecByRef(db.Source, descexec.ExecRef) (DefRec, error)
	SelectRecsByRefs(db.Source, []descexec.ExecRef) ([]DefRec, error)
	SelectRecByQN(db.Source, uniqsym.ADT) (DefRec, error)
	SelectRecsByQNs(db.Source, []uniqsym.ADT) ([]DefRec, error)
}

type defRecDS struct {
	XactID string `db:"desc_id"`
	ExpVK  int64  `db:"exp_vk"`
}
