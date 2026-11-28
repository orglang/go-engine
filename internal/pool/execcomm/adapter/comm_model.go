package adapter

import (
	"database/sql"

	commactadapter "orglang/go-engine/internal/pool/commact/adapter"
)

type CommRec struct {
	CommID         string `db:"comm_id"`
	ProviderExecID string `db:"provider_exec_id"`
	ProviderExecRN int64  `db:"provider_exec_rn"`
	ClientExecID   string `db:"client_exec_id"`
	ClientExecRN   int64  `db:"client_exec_rn"`
	OffsetNr       int64  `db:"offset_nr"`
	ActNr          int64  `db:"act_nr"`
}

type CommMod struct {
	CommID         string          `db:"comm_id"`
	OffsetNr       sql.Null[int64] `db:"offset_nr"`
	ActNr          sql.Null[int64] `db:"act_nr"`
	ProviderExecID string          `db:"-"`
	ProviderExecRN int64           `db:"-"`
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
