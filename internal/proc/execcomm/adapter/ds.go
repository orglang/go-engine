package adapter

import (
	"database/sql"

	commactadapter "orglang/go-engine/internal/proc/commact/adapter"
)

type CommMod struct {
	CommID   string          `db:"comm_id"`
	OffsetNr sql.Null[int64] `db:"offset_nr"`
}

type CommQry struct {
	CommID string           `db:"comm_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type CommSnap struct {
	CommID         string                  `db:"comm_id"`
	ProviderExecID string                  `db:"provider_exec_id"`
	ProviderExecRN int64                   `db:"provider_exec_rn"`
	OffsetNr       int64                   `db:"offset_nr"`
	ActNr          int64                   `db:"act_nr"`
	Acts           []commactadapter.ActRec `db:"acts"`
}
