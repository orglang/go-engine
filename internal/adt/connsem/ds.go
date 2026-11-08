package connsem

type SemRefDS struct {
	ConnID string `db:"conn_id"`
	ConnRN int64  `db:"conn_rn"`
}
