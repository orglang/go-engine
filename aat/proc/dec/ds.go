package dec

import (
	"orglang/orglang/avt/data"
	"orglang/orglang/avt/id"
)

type Repo interface {
	Insert(data.Source, ProcRec) error
	SelectAll(data.Source) ([]ProcRef, error)
	SelectByID(data.Source, id.ADT) (ProcSnap, error)
	SelectByIDs(data.Source, []id.ADT) ([]ProcRec, error)
	SelectEnv(data.Source, []id.ADT) (map[id.ADT]ProcRec, error)
}

type bndSpecDS struct {
	ChnlPH string `json:"chnl_ph"`
	TypeQN string `json:"role_qn"`
}

type sigRefDS struct {
	SigID string `db:"sig_id"`
	SigRN int64  `db:"rev"`
	Title string `db:"title"`
}

type sigRecDS struct {
	SigID string      `db:"sig_id"`
	Title string      `db:"title"`
	Ys    []bndSpecDS `db:"ys"`
	X     bndSpecDS   `db:"x"`
	SigRN int64       `db:"rev"`
}

type sigSnapDS struct {
	SigID string      `db:"sig_id"`
	Title string      `db:"title"`
	Ys    []bndSpecDS `db:"ys"`
	X     bndSpecDS   `db:"x"`
	SigRN int64       `db:"rev"`
}
