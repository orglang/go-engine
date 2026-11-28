package adapter

type DefRec struct {
	TypeQN string `db:"type_qn"`
	TypeID string `db:"type_id"`
	ProgID string `db:"prog_id"`
	ProgRN int64  `db:"prog_rn"`
	ExpCK  int64  `db:"expr_ck"`
}
