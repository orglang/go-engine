package adapter

import "database/sql"

type VarRec struct {
	ExecID sql.NullString `db:"exec_id"`
	ExecRN sql.NullInt64  `db:"exec_rn"`
	ConnID sql.NullString `db:"conn_id"`
	ChnlID sql.NullString `db:"chnl_id"`
	ChnlPH sql.NullString `db:"chnl_ph"`
	ExpCK  sql.NullInt64  `db:"exp_ck"`
	ChnlBS sql.NullInt16  `db:"side"`
}
