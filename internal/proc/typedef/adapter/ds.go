package adapter

type DefRec struct {
	TypeID string `db:"type_id"`
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
	ExpCK  int64  `db:"exp_ck"`
}
