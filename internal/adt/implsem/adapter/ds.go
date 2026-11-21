package adapter

type semRec struct {
	ImplQN string `db:"impl_qn"`
	ImplID string `db:"impl_id"`
	Kind   int16  `db:"kind"`
}
