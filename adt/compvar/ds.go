package compvar

import (
	"database/sql"
)

type VarRecDS struct {
	CompID sql.NullString `db:"comp_id"`
	CompRN sql.NullInt64  `db:"comp_rn"`
	CommID sql.NullString `db:"comm_id"`
	ChnlID sql.NullString `db:"chnl_id"`
	ChnlPH sql.NullString `db:"chnl_ph"`
	ExpVK  sql.NullInt64  `db:"exp_vk"`
	ChnlBS sql.NullInt16  `db:"side"`
}
