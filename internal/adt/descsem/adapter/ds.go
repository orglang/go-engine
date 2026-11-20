package adapter

type semRecDS struct {
	DescQN string `db:"desc_qn"`
	DescID string `db:"desc_id"`
	Kind   int16  `db:"kind"`
}
