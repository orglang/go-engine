package adapter

type RootRef struct {
	ExecID string `db:"exec_id"`
	ExecRN int64  `db:"exec_rn"`
}
