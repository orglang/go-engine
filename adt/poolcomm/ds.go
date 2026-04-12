package poolcomm

import (
	"database/sql"
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/uniqsym"
)

type Repo interface {
	AddRec(db.Source, ConnRec) error
	ModifyRec(db.Source, CommMod) error
	GetRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error)
	GetSnapByQry(db.Source, CommQry) (CommSnap, error)
}

type connRecDS struct {
	CommID string `db:"comm_id"`
	CommON int64  `db:"comm_on"`
}

type commModDS struct {
	CommID string          `db:"comm_id"`
	CommON sql.Null[int64] `db:"comm_on"`
}

type commQryDS struct {
	CommID string           `db:"comm_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type commSnapDS struct {
	CommID string               `db:"comm_id"`
	CommRN int64                `db:"comm_rn"`
	CommON int64                `db:"comm_on"`
	Steps  []poolstep.StepRecDS `db:"steps"`
}

type commJoinDS struct {
	CommRef commsem.SemRefDS `db:"sem"`
	CommON  int64            `db:"conn.comm_on"`
}
