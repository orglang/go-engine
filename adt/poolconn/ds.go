package poolconn

import (
	"database/sql"
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/poolstep"
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
	CommON int64  `db:"comm_on"`
}

type connModDS struct {
	CommID string          `db:"comm_id"`
	CommON sql.Null[int64] `db:"comm_on"`
}

type connQryDS struct {
	CommID string           `db:"comm_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type connSnapDS struct {
	CommID string               `db:"comm_id"`
	CommRN int64                `db:"comm_rn"`
	Steps  []poolstep.StepRecDS `db:"steps"`
}
