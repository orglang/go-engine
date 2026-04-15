package commexch

import (
	"database/sql"
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/semcomm"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/pool/commturn"
)

type Repo interface {
	AddRec(db.Source, ExchRec) error
	ModifyRec(db.Source, ExchMod) error
	GetRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]semcomm.CommRef, error)
	GetSnapByQry(db.Source, ExchQry) (ExchSnap, error)
}

type exchRecDS struct {
	CommID   string `db:"comm_id"`
	OffsetNr int64  `db:"comm_on"`
}

type exchModDS struct {
	CommID   string          `db:"comm_id"`
	OffsetNr sql.Null[int64] `db:"comm_on"`
}

type exchQryDS struct {
	CommID string           `db:"comm_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type exchSnapDS struct {
	CommID string               `db:"comm_id"`
	CommRN int64                `db:"comm_rn"`
	Turns  []commturn.TurnRecDS `db:"steps"`
}
