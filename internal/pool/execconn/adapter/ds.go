package adapter

import (
	"database/sql"

	connactadapter "orglang/go-engine/internal/pool/connact/adapter"
)

type ConnRec struct {
	ConnID         string `db:"conn_id"`
	ProviderExecID string `db:"provider_exec_id"`
	ProviderExecRN int64  `db:"provider_exec_rn"`
	ClientExecID   string `db:"client_exec_id"`
	ClientExecRN   int64  `db:"client_exec_rn"`
	OffsetNr       int64  `db:"offset_nr"`
	ActNr          int64  `db:"act_nr"`
}

type ConnMod struct {
	ConnID         string          `db:"conn_id"`
	OffsetNr       sql.Null[int64] `db:"offset_nr"`
	ActNr          sql.Null[int64] `db:"act_nr"`
	ProviderExecID string          `db:"-"`
	ProviderExecRN int64           `db:"-"`
}

type ConnQry struct {
	ConnID string           `db:"conn_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type ConnSnap struct {
	ConnID         string                  `db:"conn_id"`
	ProviderExecID string                  `db:"provider_exec_id"`
	ProviderExecRN int64                   `db:"provider_exec_rn"`
	OffsetNr       int64                   `db:"offset_nr"`
	ActNr          int64                   `db:"act_nr"`
	Acts           []connactadapter.ActRec `db:"acts"`
}
