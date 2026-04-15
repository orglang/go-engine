package compvar

import (
	"database/sql"
)

type VarRecDS struct {
	ImplID sql.NullString `db:"impl_id"`
	ImplRN sql.NullInt64  `db:"impl_rn"`
	CommID sql.NullString `db:"comm_id"`
	ChnlID sql.NullString `db:"chnl_id"`
	ChnlPH sql.NullString `db:"chnl_ph"`
	ExpVK  sql.NullInt64  `db:"exp_vk"`
	ChnlBS sql.NullInt16  `db:"side"`
}
