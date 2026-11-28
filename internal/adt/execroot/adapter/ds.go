package adapter

type RootRef struct {
	RootID string `db:"exec_id"`
	RootRN int64  `db:"exec_rn"`
}
