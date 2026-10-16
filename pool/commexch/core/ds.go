package core

import (
	"database/sql"
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/uniqsym"
	commturn "orglang/go-engine/pool/commturn/core"
)

type Repo interface {
	AddRec(db.UoW, ExchRec) error
	ModifyRec(db.UoW, ExchMod) error
	GetRefsByQNs(db.UoW, []uniqsym.ADT) (map[uniqsym.ADT]commsem.SemRef, error)
	GetSnapByQry(db.UoW, ExchQry) (ExchSnap, error)
}

type ExchRecDS struct {
	CommID   string `db:"comm_id"`
	CommRN   int64  `db:"comm_rn"`
	OffsetNr int64  `db:"offset_nr"`
}

type ExchModDS struct {
	CommID   string          `db:"comm_id"`
	OffsetNr sql.Null[int64] `db:"offset_nr"`
}

type ExchQryDS struct {
	CommID string           `db:"comm_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type ExchSnapDS struct {
	CommID   string               `db:"comm_id"`
	CommRN   int64                `db:"comm_rn"`
	OffsetNr int64                `db:"offset_nr"`
	Turns    []commturn.TurnRecDS `db:"turns"`
}
