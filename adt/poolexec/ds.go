package poolexec

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	InsertRec(db.Source, ExecRec) error
	TouchRec(db.Source, implsem.SemRef) error
	TouchRecs(db.Source, []implsem.SemRef) error
	SelectRecsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]ExecRec, error)
	SelectRefs(db.Source) ([]implsem.SemRef, error)
	SelectSubs(db.Source, implsem.SemRef) (ExecSnap, error)
}

// TODO очень напоминает implvar.VarRecDS
type execRecDS struct {
	ChnlID string `db:"chnl_id"`
	ChnlPH string `db:"chnl_ph"`
	ExpVK  int64  `db:"exp_vk"`
	ImplID string `db:"impl_id"`
	ImplRN int64  `db:"impl_rn"`
}

type execSnapDS struct {
	ImplID   string             `db:"impl_id"`
	ImplRN   int64              `db:"impl_rn"`
	Title    string             `db:"title"`
	SubExecs []implsem.SemRefDS `db:"subs"`
}
