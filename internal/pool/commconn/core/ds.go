package core

import (
	"context"
	"database/sql"

	"orglang/go-engine/internal/lib/tx"

	commact "orglang/go-engine/internal/pool/commact/core"
)

type Repo interface {
	AddRec(context.Context, tx.Conn, ConnRec) error
	ModifyRec(context.Context, tx.Conn, ConnMod) error
	GetSnapByQry(context.Context, tx.Conn, ConnQry) (ConnSnap, error)
}

type ConnRecDS struct {
	ConnID         string `db:"conn_id"`
	ProviderCompID string `db:"provider_comp_id"`
	ProviderCompRN int64  `db:"provider_comp_rn"`
	ClientCompID   string `db:"client_comp_id"`
	ClientCompRN   int64  `db:"client_comp_rn"`
	OffsetNr       int64  `db:"offset_nr"`
	ActNr          int64  `db:"act_nr"`
}

type ConnModDS struct {
	ConnID         string          `db:"conn_id"`
	OffsetNr       sql.Null[int64] `db:"offset_nr"`
	ActNr          sql.Null[int64] `db:"act_nr"`
	ProviderCompID string          `db:"-"`
	ProviderCompRN int64           `db:"-"`
}

type ConnQryDS struct {
	ConnID string           `db:"conn_id"`
	ChnlID sql.Null[string] `db:"chnl_id"`
}

type ConnSnapDS struct {
	ConnID         string               `db:"conn_id"`
	ProviderCompID string               `db:"provider_comp_id"`
	ProviderCompRN int64                `db:"provider_comp_rn"`
	OffsetNr       int64                `db:"offset_nr"`
	ActNr          int64                `db:"act_nr"`
	Acts          []commact.ActRecDS `db:"turns"`
}
