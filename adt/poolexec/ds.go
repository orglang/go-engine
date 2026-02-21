package poolexec

import (
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/lib/db"
)

// Port
type Repo interface {
	InsertRec(db.Source, ExecRec) error
	InsertLiab(db.Source, Liab) error
	SelectRefs(db.Source) ([]implsem.SemRef, error)
	SelectSubs(db.Source, implsem.SemRef) (ExecSnap, error)
}

type execRefDS = struct {
	ID string `db:"exec_id"`
	RN int64  `db:"exec_rn"`
}

type execSnapDS struct {
	ImplID   string             `db:"exec_id"`
	ImplRN   int64              `db:"exec_rn"`
	Title    string             `db:"title"`
	SubExecs []implsem.SemRefDS `db:"subs"`
}

type execRecDS struct {
	ImplID string `db:"exec_id"`
	ImplRN int64  `db:"exec_rn"`
}

type liabDS struct {
	ImplID string `db:"exec_id"`
	ImplRN int64  `db:"exec_rn"`
	ProcID string `db:"proc_id"`
}
