package adapter

import (
	"database/sql"

	connactadapter "orglang/go-engine/internal/proc/connact/adapter"
)

type ConnMod struct {
	ConnID   string          `db:"conn_id"`
	OffsetNr sql.Null[int64] `db:"offset_nr"`
}

type ConnQry struct {
	ConnID string           `db:"conn_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type ConnSnap struct {
	ConnID         string                  `db:"conn_id"`
	ProviderCompID string                  `db:"provider_comp_id"`
	ProviderCompRN int64                   `db:"provider_comp_rn"`
	OffsetNr       int64                   `db:"offset_nr"`
	ActNr          int64                   `db:"act_nr"`
	Acts           []connactadapter.ActRec `db:"acts"`
}
