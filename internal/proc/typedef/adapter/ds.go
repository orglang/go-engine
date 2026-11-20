package adapter

type DefRecDS struct {
	TypeID string `db:"type_id"`
	TypeRN int64  `db:"type_rn"`
	ExpCK  int64  `db:"exp_ck"`
}
